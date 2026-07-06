package gen

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tamnd/tama/pkg/config"
)

func testCfg(baseURL string) config.LLM {
	return config.LLM{
		BaseURL:        baseURL,
		APIKey:         "test-key",
		Model:          "test-model",
		RequestTimeout: 5 * time.Second,
		ConnectTimeout: time.Second,
	}
}

func newTestClient(cfg config.LLM) *Client {
	c := New(cfg)
	c.backoffBase = time.Millisecond
	return c
}

func TestPingSendsAuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if r.Method != http.MethodGet || r.URL.Path != "/v1/models" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		w.Write([]byte(`{"data":[{"id":"model-a"},{"id":"model-b"}]}`))
	}))
	defer srv.Close()

	models, err := newTestClient(testCfg(srv.URL+"/v1")).Ping(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer test-key" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	if len(models) != 2 || models[0].ID != "model-a" {
		t.Errorf("models = %+v", models)
	}
}

func TestNoAuthHeaderWithoutKey(t *testing.T) {
	var sawAuth bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, sawAuth = r.Header["Authorization"]
		w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	cfg := testCfg(srv.URL)
	cfg.APIKey = ""
	if _, err := newTestClient(cfg).Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
	if sawAuth {
		t.Error("Authorization header sent despite empty key")
	}
}

func TestRetryOn429HonorsRetryAfter(t *testing.T) {
	var calls int
	start := time.Now()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Write([]byte(`{"data":[{"id":"m"}]}`))
	}))
	defer srv.Close()

	models, err := newTestClient(testCfg(srv.URL)).Ping(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if calls != 3 || len(models) != 1 {
		t.Errorf("calls = %d, models = %+v", calls, models)
	}
	// Two Retry-After: 1 waits must actually have happened.
	if elapsed := time.Since(start); elapsed < 2*time.Second {
		t.Errorf("finished in %v, Retry-After ignored", elapsed)
	}
}

func TestGivesUpAfterThreeAttempts(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	_, err := newTestClient(testCfg(srv.URL)).Ping(context.Background())
	if err == nil || !strings.Contains(err.Error(), "gave up after 3 attempts") {
		t.Fatalf("err = %v", err)
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	_, err := newTestClient(testCfg(srv.URL)).Ping(context.Background())
	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("err = %v", err)
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1 (4xx must not retry)", calls)
	}
}

func TestTimeoutSurfaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer srv.Close()

	cfg := testCfg(srv.URL)
	cfg.RequestTimeout = 100 * time.Millisecond
	cfg.ConnectTimeout = 50 * time.Millisecond
	_, err := newTestClient(cfg).Ping(context.Background())
	if err == nil {
		t.Fatal("slow endpoint did not error")
	}
}

func TestCompleteJSONModeBody(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("content type = %q", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Error(err)
		}
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"{\"ok\":true}"}}]}`))
	}))
	defer srv.Close()

	out, err := newTestClient(testCfg(srv.URL)).Complete(context.Background(), CompleteRequest{
		Messages:    []Message{{Role: "user", Content: "make a lesson"}},
		Temperature: 0.7,
		MaxTokens:   512,
		JSONMode:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"ok":true}` {
		t.Errorf("content = %q", out)
	}

	if got["model"] != "test-model" || got["temperature"] != 0.7 || got["max_tokens"] != float64(512) {
		t.Errorf("body = %+v", got)
	}
	rf, ok := got["response_format"].(map[string]any)
	if !ok || rf["type"] != "json_object" {
		t.Errorf("response_format = %+v", got["response_format"])
	}
}

func TestCompleteWithoutModelFails(t *testing.T) {
	cfg := testCfg("http://127.0.0.1:1")
	cfg.Model = ""
	_, err := newTestClient(cfg).Complete(context.Background(), CompleteRequest{})
	if err == nil || !strings.Contains(err.Error(), "no model configured") {
		t.Fatalf("err = %v", err)
	}
}
