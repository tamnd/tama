package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteDataEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	WriteData(w, map[string]string{"status": "ok"})
	if w.Code != http.StatusOK {
		t.Errorf("status = %d", w.Code)
	}
	if got := strings.TrimSpace(w.Body.String()); got != `{"data":{"status":"ok"}}` {
		t.Errorf("body = %s", got)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("content type = %q", ct)
	}
}

func TestWritePageCursor(t *testing.T) {
	w := httptest.NewRecorder()
	WritePage(w, []int{1, 2}, "abc")
	if !strings.Contains(w.Body.String(), `"next_cursor":"abc"`) {
		t.Errorf("body = %s", w.Body.String())
	}

	w = httptest.NewRecorder()
	WritePage(w, []int{1, 2}, "")
	if strings.Contains(w.Body.String(), "next_cursor") {
		t.Errorf("last page still carries next_cursor: %s", w.Body.String())
	}
}

func TestWriteErrorStatusMatchesCode(t *testing.T) {
	want := map[Code]int{
		CodeBadRequest:   400,
		CodeUnauthorized: 401,
		CodeForbidden:    403,
		CodeNotFound:     404,
		CodeConflict:     409,
		CodeRateLimited:  429,
		CodeInternal:     500,
	}
	for code, status := range want {
		w := httptest.NewRecorder()
		WriteError(w, code, "boom")
		if w.Code != status {
			t.Errorf("%s status = %d, want %d", code, w.Code, status)
		}
		if !strings.Contains(w.Body.String(), string(code)) {
			t.Errorf("%s body = %s", code, w.Body.String())
		}
	}
}

func TestReadJSONRejectsGarbage(t *testing.T) {
	for name, body := range map[string]string{
		"malformed": `{"username":`,
		"trailing":  `{"username":"a"}{"more":true}`,
		"unknown":   `{"nope":true}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		var dst struct {
			Username string `json:"username"`
		}
		if err := ReadJSON(httptest.NewRecorder(), req, &dst); err == nil {
			t.Errorf("%s: ReadJSON accepted %q", name, body)
		}
	}
}

func TestReadJSONAcceptsCharsetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username":"a"}`))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	var dst struct {
		Username string `json:"username"`
	}
	if err := ReadJSON(httptest.NewRecorder(), req, &dst); err != nil {
		t.Fatalf("ReadJSON: %v", err)
	}
}

func TestCursorRoundTrip(t *testing.T) {
	c := EncodeCursor("user:42")
	got, err := DecodeCursor(c)
	if err != nil || got != "user:42" {
		t.Errorf("round trip = %q, %v", got, err)
	}
	if got, err := DecodeCursor(""); err != nil || got != "" {
		t.Errorf("empty cursor = %q, %v", got, err)
	}
	if _, err := DecodeCursor("!!!not-base64!!!"); err == nil {
		t.Error("garbled cursor accepted")
	}
}

func TestParseLimit(t *testing.T) {
	cases := []struct {
		query string
		want  int
		bad   bool
	}{
		{"", 20, false},
		{"limit=5", 5, false},
		{"limit=100", 100, false},
		{"limit=500", 100, false},
		{"limit=0", 0, true},
		{"limit=puppy", 0, true},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "/?"+tc.query, nil)
		got, err := ParseLimit(req)
		if tc.bad {
			if err == nil {
				t.Errorf("%q: no error", tc.query)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Errorf("%q = %d, %v, want %d", tc.query, got, err, tc.want)
		}
	}
}

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(hash, "$argon2id$v=19$m=65536,t=3,p=2$") {
		t.Errorf("hash = %q, want PHC argon2id prefix with the doc parameters", hash)
	}
	if !VerifyPassword(hash, "correct horse") {
		t.Error("correct password rejected")
	}
	if VerifyPassword(hash, "wrong horse") {
		t.Error("wrong password accepted")
	}
	if VerifyPassword("$argon2id$v=19$m=65536,t=3,p=2$garbage", "correct horse") {
		t.Error("truncated hash accepted")
	}
	if VerifyPassword("", "correct horse") {
		t.Error("empty hash accepted")
	}
}

func TestHashesAreSalted(t *testing.T) {
	a, _ := HashPassword("same password")
	b, _ := HashPassword("same password")
	if a == b {
		t.Error("two hashes of one password match; salt is not random")
	}
}
