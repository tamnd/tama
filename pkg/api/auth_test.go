package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tamnd/tama/pkg/store"
)

// newTestAPI wires the real mux, middleware, and a real temp-dir store, the
// same stack serve runs minus pkg/server's outer middleware.
func newTestAPI(t *testing.T) (*Handlers, http.Handler) {
	t.Helper()
	db, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "tama.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	h := New(db)
	mux := http.NewServeMux()
	h.Routes(mux)
	return h, h.SessionMiddleware(mux)
}

func postJSON(t *testing.T, handler http.Handler, path, body string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func decodeError(t *testing.T, w *httptest.ResponseRecorder) (code, message string) {
	t.Helper()
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("error envelope: %v in %q", err, w.Body.String())
	}
	return body.Error.Code, body.Error.Message
}

func sessionFrom(t *testing.T, w *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range w.Result().Cookies() {
		if c.Name == SessionCookie && c.Value != "" {
			return c
		}
	}
	t.Fatalf("no %s cookie in response", SessionCookie)
	return nil
}

func TestRegisterHappyPath(t *testing.T) {
	_, handler := newTestAPI(t)

	w := postJSON(t, handler, "/api/auth/register", `{"username":"Mochi","password":"password1"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("register = %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Data UserPayload `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Data.Username != "mochi" {
		t.Errorf("username = %q, want lowercased mochi", body.Data.Username)
	}
	if strings.Contains(w.Body.String(), "password") {
		t.Error("response leaks a password field")
	}

	c := sessionFrom(t, w)
	if !c.HttpOnly || c.SameSite != http.SameSiteLaxMode || c.Path != "/" {
		t.Errorf("cookie attributes = %+v", c)
	}
	if c.Secure {
		t.Error("cookie is Secure over plain HTTP")
	}
}

func TestRegisterValidation(t *testing.T) {
	_, handler := newTestAPI(t)

	cases := []struct {
		name, body, wantIn string
	}{
		{"short username", `{"username":"ab","password":"password1"}`, "username"},
		{"bad chars", `{"username":"mo chi!","password":"password1"}`, "username"},
		{"short password", `{"username":"mochi","password":"short"}`, "password"},
	}
	for _, tc := range cases {
		w := postJSON(t, handler, "/api/auth/register", tc.body)
		if w.Code != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want 400", tc.name, w.Code)
		}
		code, msg := decodeError(t, w)
		if code != "bad_request" || !strings.Contains(msg, tc.wantIn) {
			t.Errorf("%s: error = %s %q", tc.name, code, msg)
		}
	}
}

func TestRegisterDuplicateConflicts(t *testing.T) {
	_, handler := newTestAPI(t)

	postJSON(t, handler, "/api/auth/register", `{"username":"mochi","password":"password1"}`)
	w := postJSON(t, handler, "/api/auth/register", `{"username":"MOCHI","password":"password2"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("duplicate register = %d, want 409", w.Code)
	}
	if code, _ := decodeError(t, w); code != "conflict" {
		t.Errorf("code = %s, want conflict", code)
	}
}

func TestRegisterRejectsNonJSON(t *testing.T) {
	_, handler := newTestAPI(t)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader("username=mochi"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("form post = %d, want 415", w.Code)
	}
	if code, _ := decodeError(t, w); code != "bad_request" {
		t.Errorf("code = %s, want bad_request", code)
	}
}

func TestRegisterBodySizeCap(t *testing.T) {
	_, handler := newTestAPI(t)

	huge := fmt.Sprintf(`{"username":"mochi","password":"password1","x":%q}`, strings.Repeat("a", maxBodyBytes))
	w := postJSON(t, handler, "/api/auth/register", huge)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("oversized body = %d, want 400", w.Code)
	}
	if _, msg := decodeError(t, w); !strings.Contains(msg, "exceeds") {
		t.Errorf("message = %q, want size cap mention", msg)
	}
}

func TestLoginAndMeAndLogout(t *testing.T) {
	_, handler := newTestAPI(t)
	postJSON(t, handler, "/api/auth/register", `{"username":"mochi","password":"password1"}`)

	w := postJSON(t, handler, "/api/auth/login", `{"username":"mochi","password":"password1"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("login = %d: %s", w.Code, w.Body.String())
	}
	cookie := sessionFrom(t, w)

	// me with the cookie
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"username":"mochi"`) {
		t.Fatalf("me = %d: %s", rec.Code, rec.Body.String())
	}

	// logout kills the session
	w = postJSON(t, handler, "/api/auth/logout", ``, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("logout = %d: %s", w.Code, w.Body.String())
	}
	req = httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("me after logout = %d, want 401", rec.Code)
	}
}

func TestLoginBadCredentials(t *testing.T) {
	_, handler := newTestAPI(t)
	postJSON(t, handler, "/api/auth/register", `{"username":"mochi","password":"password1"}`)

	for name, body := range map[string]string{
		"wrong password": `{"username":"mochi","password":"password2"}`,
		"no such user":   `{"username":"nobody","password":"password1"}`,
	} {
		w := postJSON(t, handler, "/api/auth/login", body)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s: status = %d, want 401", name, w.Code)
		}
		// The message must not say which half was wrong.
		if _, msg := decodeError(t, w); msg != "invalid credentials" {
			t.Errorf("%s: message = %q", name, msg)
		}
	}
}

func TestLoginRotatesSession(t *testing.T) {
	_, handler := newTestAPI(t)
	postJSON(t, handler, "/api/auth/register", `{"username":"mochi","password":"password1"}`)

	first := sessionFrom(t, postJSON(t, handler, "/api/auth/login", `{"username":"mochi","password":"password1"}`))
	second := sessionFrom(t, postJSON(t, handler, "/api/auth/login", `{"username":"mochi","password":"password1"}`, first))
	if first.Value == second.Value {
		t.Fatal("login reused the session token")
	}

	// The old token must be dead, not just superseded.
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(first)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("old session after rotation = %d, want 401", rec.Code)
	}
}

func TestMeWithoutSession(t *testing.T) {
	_, handler := newTestAPI(t)
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("me = %d, want 401", w.Code)
	}
	if code, _ := decodeError(t, w); code != "unauthorized" {
		t.Errorf("code = %s, want unauthorized", code)
	}
}

func TestLogoutWithoutSession(t *testing.T) {
	_, handler := newTestAPI(t)
	w := postJSON(t, handler, "/api/auth/logout", ``)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("logout without session = %d, want 401", w.Code)
	}
}

func TestAuthRateLimit(t *testing.T) {
	_, handler := newTestAPI(t)

	// Freeze the clock so the bucket cannot refill mid-test, and post empty
	// bodies: they fail fast after the limiter, no hashing involved.
	now := store.Now()
	old := store.Now
	store.Now = func() time.Time { return now }
	t.Cleanup(func() { store.Now = old })

	var last *httptest.ResponseRecorder
	for i := 0; i < 11; i++ {
		last = postJSON(t, handler, "/api/auth/login", ``)
	}
	if last.Code != http.StatusTooManyRequests {
		t.Fatalf("11th attempt = %d, want 429", last.Code)
	}
	if code, _ := decodeError(t, last); code != "rate_limited" {
		t.Errorf("code = %s, want rate_limited", code)
	}
	if last.Header().Get("Retry-After") == "" {
		t.Error("429 without Retry-After header")
	}

	// Register shares the same bucket, so it is limited too now.
	w := postJSON(t, handler, "/api/auth/register", `{"username":"mochi","password":"password1"}`)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("register on drained bucket = %d, want 429", w.Code)
	}

	// A different IP still has a full bucket.
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(``))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.7:9999"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code == http.StatusTooManyRequests {
		t.Errorf("fresh IP = %d, want anything but 429", rec.Code)
	}
}
