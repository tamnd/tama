package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tamnd/tama/pkg/store"
)

// SessionCookie is the session cookie name. Its value is sessions.token.
const SessionCookie = "tama_session"

// sessionTTL is how long a session lives, 30 days.
const sessionTTL = 30 * 24 * time.Hour

// Handlers owns the API endpoints and their shared state.
type Handlers struct {
	db      *store.DB
	authRL  *rateLimiter
	Version string
}

// New builds the handler set over the store.
func New(db *store.DB) *Handlers {
	return &Handlers{db: db, authRL: newRateLimiter(10)}
}

// Routes mounts every /api route onto mux. Session resolution happens in the
// middleware pkg/server wraps around the whole chain; routes that need a
// user call currentUser themselves.
func (h *Handlers) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/register", h.authRL.limit(h.handleRegister))
	mux.HandleFunc("POST /api/auth/login", h.authRL.limit(h.handleLogin))
	mux.HandleFunc("POST /api/auth/logout", h.requireUser(h.handleLogout))
	mux.HandleFunc("GET /api/me", h.requireUser(h.handleMe))
}

// UserPayload is the client-facing user object; the hash never leaves the
// server.
type UserPayload struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
	IsAdmin     bool   `json:"is_admin"`
	CreatedAt   int64  `json:"created_at"`
}

func userPayload(u store.User) UserPayload {
	return UserPayload{
		ID: u.ID, Username: u.Username, DisplayName: u.DisplayName,
		Avatar: u.Avatar, IsAdmin: u.IsAdmin, CreatedAt: u.CreatedAt,
	}
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var usernameRe = regexp.MustCompile(`^[a-z0-9_]{3,24}$`)

// ValidateUsername lowercases and checks a username, returning the canonical
// form. The CLI user commands share it with the register endpoint.
func ValidateUsername(raw string) (string, error) {
	name := strings.ToLower(strings.TrimSpace(raw))
	if !usernameRe.MatchString(name) {
		return "", errf(CodeBadRequest, "username must be 3-24 characters of a-z, 0-9, or _")
	}
	return name, nil
}

// ValidatePassword enforces the minimum length.
func ValidatePassword(pw string) error {
	if len(pw) < 8 {
		return errf(CodeBadRequest, "password must be at least 8 characters")
	}
	return nil
}

func (h *Handlers) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req credentials
	if err := ReadJSON(w, r, &req); err != nil {
		writeErr(w, r, err)
		return
	}
	username, err := ValidateUsername(req.Username)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	if err := ValidatePassword(req.Password); err != nil {
		writeErr(w, r, err)
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		writeErr(w, r, internalErr(err))
		return
	}
	user, err := h.db.CreateUser(r.Context(), username, hash, false)
	if err != nil {
		if isUniqueViolation(err) {
			writeErr(w, r, errf(CodeConflict, "username %s is taken", username))
		} else {
			writeErr(w, r, internalErr(err))
		}
		return
	}

	if err := h.startSession(w, r, user.ID); err != nil {
		writeErr(w, r, internalErr(err))
		return
	}
	slog.InfoContext(r.Context(), "user registered", "username", username)
	WriteData(w, userPayload(user))
}

func (h *Handlers) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req credentials
	if err := ReadJSON(w, r, &req); err != nil {
		writeErr(w, r, err)
		return
	}

	// One message for both bad username and bad password, and a hash check
	// either way so the two cases take the same time.
	badCreds := errf(CodeUnauthorized, "invalid credentials")
	username := strings.ToLower(strings.TrimSpace(req.Username))
	user, err := h.db.UserByUsername(r.Context(), username)
	if errors.Is(err, store.ErrNotFound) {
		VerifyPassword(dummyHash, req.Password)
		writeErr(w, r, badCreds)
		return
	}
	if err != nil {
		writeErr(w, r, internalErr(err))
		return
	}
	if !VerifyPassword(user.PasswordHash, req.Password) {
		writeErr(w, r, badCreds)
		return
	}

	// A fresh token on every login; the old cookie, if any, dies with its
	// session row.
	if c, err := r.Cookie(SessionCookie); err == nil {
		h.db.DeleteSession(r.Context(), c.Value)
	}
	if err := h.startSession(w, r, user.ID); err != nil {
		writeErr(w, r, internalErr(err))
		return
	}
	slog.InfoContext(r.Context(), "user logged in", "username", username)
	WriteData(w, userPayload(user))
}

func (h *Handlers) handleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(SessionCookie); err == nil {
		if err := h.db.DeleteSession(r.Context(), c.Value); err != nil {
			writeErr(w, r, internalErr(err))
			return
		}
	}
	http.SetCookie(w, sessionCookie(r, "", -1))
	WriteData(w, map[string]bool{"ok": true})
}

func (h *Handlers) handleMe(w http.ResponseWriter, r *http.Request) {
	user, _ := CurrentUser(r.Context())
	WriteData(w, userPayload(user))
}

func (h *Handlers) startSession(w http.ResponseWriter, r *http.Request, userID int64) error {
	sess, err := h.db.CreateSession(r.Context(), userID, r.UserAgent(), sessionTTL)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessionCookie(r, sess.Token, int(sessionTTL.Seconds())))
	return nil
}

// sessionCookie builds the tama_session cookie; Secure flips on when the
// request itself came over TLS.
func sessionCookie(r *http.Request, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookie,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	}
}

type ctxKey int

const userKey ctxKey = 0

// SessionMiddleware resolves the session cookie to a user and stashes it in
// the request context. It never rejects; requireUser does that per route.
func (h *Handlers) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(SessionCookie); err == nil {
			if _, user, err := h.db.SessionByToken(r.Context(), c.Value); err == nil {
				r = r.WithContext(context.WithValue(r.Context(), userKey, user))
			}
		}
		next.ServeHTTP(w, r)
	})
}

// CurrentUser pulls the authenticated user out of the context.
func CurrentUser(ctx context.Context) (store.User, bool) {
	u, ok := ctx.Value(userKey).(store.User)
	return u, ok
}

// requireUser guards a route: 401 unless the middleware resolved a session.
func (h *Handlers) requireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := CurrentUser(r.Context()); !ok {
			WriteError(w, CodeUnauthorized, "authentication required")
			return
		}
		next(w, r)
	}
}

// RequireUser is requireUser for routes other packages mount.
func (h *Handlers) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return h.requireUser(next)
}

// isUniqueViolation spots SQLite's UNIQUE constraint error without importing
// the driver's error types here.
func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
