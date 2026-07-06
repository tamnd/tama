// End-to-end: boot the real server on port 0 with a temp data dir and drive
// the auth flow plus the SPA index over real HTTP, cookies and all.
package tama_test

import (
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tamnd/tama/pkg/config"
	"github.com/tamnd/tama/pkg/server"
	"github.com/tamnd/tama/pkg/store"
)

func TestEndToEnd(t *testing.T) {
	dir := t.TempDir()
	db, err := store.Open(context.Background(), filepath.Join(dir, "tama.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	cfg := &config.Config{Addr: "127.0.0.1:0", DataDir: dir}
	srv := server.New(cfg, db, server.Options{Version: "e2e"})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- srv.Run(ctx) }()
	defer func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	}()

	select {
	case <-srv.Ready():
	case <-time.After(2 * time.Second):
		t.Fatal("server not ready within 2s")
	}
	base := "http://" + srv.BoundAddr()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Jar: jar}

	post := func(path, body string) *http.Response {
		t.Helper()
		resp, err := client.Post(base+path, "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("POST %s: %v", path, err)
		}
		return resp
	}
	read := func(resp *http.Response) string {
		t.Helper()
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	// Register sets the session cookie.
	resp := post("/api/auth/register", `{"username":"mochi","password":"password1"}`)
	if body := read(resp); resp.StatusCode != http.StatusOK || !strings.Contains(body, `"username":"mochi"`) {
		t.Fatalf("register = %d: %s", resp.StatusCode, body)
	}

	// me works off the cookie the jar kept.
	resp, err = client.Get(base + "/api/me")
	if err != nil {
		t.Fatal(err)
	}
	if body := read(resp); resp.StatusCode != http.StatusOK || !strings.Contains(body, `"username":"mochi"`) {
		t.Fatalf("me = %d: %s", resp.StatusCode, body)
	}

	// Fresh login rotates the session and keeps working.
	resp = post("/api/auth/login", `{"username":"mochi","password":"password1"}`)
	if body := read(resp); resp.StatusCode != http.StatusOK {
		t.Fatalf("login = %d: %s", resp.StatusCode, body)
	}

	// Logout kills the session; me turns 401.
	resp = post("/api/auth/logout", ``)
	if body := read(resp); resp.StatusCode != http.StatusOK {
		t.Fatalf("logout = %d: %s", resp.StatusCode, body)
	}
	resp, err = client.Get(base + "/api/me")
	if err != nil {
		t.Fatal(err)
	}
	if body := read(resp); resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("me after logout = %d, want 401: %s", resp.StatusCode, body)
	}

	// The SPA index serves over the same listener.
	resp, err = client.Get(base + "/")
	if err != nil {
		t.Fatal(err)
	}
	if body := read(resp); resp.StatusCode != http.StatusOK || !strings.Contains(strings.ToLower(body), "<!doctype html") {
		t.Fatalf("SPA index = %d: %.120s", resp.StatusCode, body)
	}
}
