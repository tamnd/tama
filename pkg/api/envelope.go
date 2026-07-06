package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"mime"
	"net/http"
	"strconv"
)

// maxBodyBytes caps JSON request bodies at 1 MiB.
const maxBodyBytes = 1 << 20

// WriteData writes a {"data": v} envelope with status 200.
func WriteData(w http.ResponseWriter, v any) {
	writeJSON(w, http.StatusOK, map[string]any{"data": v})
}

// WritePage writes a list envelope; nextCursor is omitted on the last page.
func WritePage(w http.ResponseWriter, v any, nextCursor string) {
	body := map[string]any{"data": v}
	if nextCursor != "" {
		body["next_cursor"] = nextCursor
	}
	writeJSON(w, http.StatusOK, body)
}

// WriteError writes the {"error":{code,message}} envelope with the status
// the code implies.
func WriteError(w http.ResponseWriter, code Code, msg string) {
	writeErrorStatus(w, code.status(), code, msg)
}

func writeErrorStatus(w http.ResponseWriter, status int, code Code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": string(code), "message": msg}})
}

// writeErr is the one place an error becomes a response. Wrapped causes log
// server-side and never render to clients.
func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	var ae *apiError
	if !errors.As(err, &ae) {
		ae = internalErr(err)
	}
	if ae.Err != nil || ae.Code == CodeInternal {
		slog.ErrorContext(r.Context(), "request failed", "code", ae.Code, "err", ae.Err)
	}
	status := ae.Status
	if status == 0 {
		status = ae.Code.status()
	}
	writeErrorStatus(w, status, ae.Code, ae.Message)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// ReadJSON decodes the request body into dst. It enforces the 1 MiB body
// cap, requires an application/json content type (415 otherwise), and
// rejects unknown fields and trailing garbage.
func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	ct := r.Header.Get("Content-Type")
	if mt, _, err := mime.ParseMediaType(ct); err != nil || mt != "application/json" {
		return &apiError{Code: CodeBadRequest, Status: http.StatusUnsupportedMediaType,
			Message: "content type must be application/json"}
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			return errf(CodeBadRequest, "request body exceeds %d bytes", tooBig.Limit)
		}
		return errf(CodeBadRequest, "malformed JSON body")
	}
	if dec.More() {
		return errf(CodeBadRequest, "request body must be a single JSON object")
	}
	return nil
}

// Cursor helpers: cursors are opaque base64 of the last row's sort key.

// EncodeCursor turns a sort key into an opaque cursor.
func EncodeCursor(key string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(key))
}

// DecodeCursor recovers the sort key; a garbled cursor is a bad_request.
func DecodeCursor(cursor string) (string, error) {
	if cursor == "" {
		return "", nil
	}
	b, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return "", errf(CodeBadRequest, "invalid cursor")
	}
	return string(b), nil
}

// ParseLimit reads ?limit= with a default of 20, capped at 100.
func ParseLimit(r *http.Request) (int, error) {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return 20, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 0, errf(CodeBadRequest, "limit must be a positive integer")
	}
	return min(n, 100), nil
}
