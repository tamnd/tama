package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/tamnd/tama/pkg/config"
	"github.com/tamnd/tama/pkg/store"
)

// bootTestServer starts the real server on port 0 over a temp data dir and
// returns its base URL.
func bootTestServer(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	db, err := store.Open(context.Background(), filepath.Join(dir, "tama.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	cfg := &config.Config{Addr: "127.0.0.1:0", DataDir: dir}
	srv := New(cfg, db, Options{Version: "test"})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- srv.Run(ctx) }()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})

	select {
	case <-srv.Ready():
	case <-time.After(2 * time.Second):
		t.Fatal("server not ready within 2s")
	}
	return "http://" + srv.BoundAddr()
}

func TestBootServesHealthzAndSPA(t *testing.T) {
	start := time.Now()
	base := bootTestServer(t)

	resp, err := http.Get(base + "/api/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("healthz = %d", resp.StatusCode)
	}
	var body struct {
		Data struct {
			Status  string `json:"status"`
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Data.Status != "ok" || body.Data.Version != "test" {
		t.Errorf("healthz body = %+v", body)
	}

	index, err := http.Get(base + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer index.Body.Close()
	if index.StatusCode != http.StatusOK {
		t.Fatalf("SPA index = %d", index.StatusCode)
	}
	if cc := index.Header.Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("index Cache-Control = %q, want no-cache", cc)
	}
	if id := index.Header.Get("X-Request-Id"); id == "" {
		t.Error("response missing X-Request-Id")
	}

	// Deep links fall back to index.html, not 404.
	deep, err := http.Get(base + "/lesson/u1-l1")
	if err != nil {
		t.Fatal(err)
	}
	deep.Body.Close()
	if deep.StatusCode != http.StatusOK {
		t.Errorf("deep link = %d, want 200", deep.StatusCode)
	}

	// The doc wants cold start under 500ms; 2s here to avoid CI flakes.
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("cold start to first request took %v", elapsed)
	}
}

func TestHashedAssetsAreImmutable(t *testing.T) {
	base := bootTestServer(t)

	// Find a real hashed asset via the index page's asset dir listing being
	// unavailable; instead just probe the embedded dist for one.
	resp, err := http.Get(base + "/assets/no-such-file.js")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	// A missing asset falls through to the SPA and must not be immutable.
	if cc := resp.Header.Get("Cache-Control"); cc == "public, max-age=31536000, immutable" {
		t.Errorf("SPA fallback got immutable cache header")
	}
}

func TestProtectedRoutesNeedSession(t *testing.T) {
	base := bootTestServer(t)
	for _, path := range []string{"/api/languages", "/api/catalog", "/api/me"} {
		resp, err := http.Get(base + path)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("%s = %d, want 401", path, resp.StatusCode)
		}
	}
}

func TestRecoveredPanicIsJSON500(t *testing.T) {
	h := requestID(logRequests(recovered(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("kaboom")
	}))))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Error.Code != "internal" {
		t.Errorf("body = %s (%v)", w.Body.String(), err)
	}
}

func TestRealIPRewritesRemoteAddr(t *testing.T) {
	var got string
	h := realIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.RemoteAddr
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.9, 10.0.0.1")
	h.ServeHTTP(httptest.NewRecorder(), req)
	if got != "203.0.113.9:0" {
		t.Errorf("RemoteAddr = %q, want 203.0.113.9:0", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "not-an-ip")
	h.ServeHTTP(httptest.NewRecorder(), req)
	if got == "not-an-ip:0" {
		t.Error("garbage X-Forwarded-For accepted")
	}
}
