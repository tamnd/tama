package store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"
)

// Session mirrors the sessions table.
type Session struct {
	Token     string
	UserID    int64
	CreatedAt int64
	ExpiresAt int64
	UserAgent string
}

// newSessionToken returns 32 random bytes as unpadded base64url.
func newSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

const createSessionSQL = `INSERT INTO sessions (token, user_id, created_at, expires_at, user_agent) VALUES (?, ?, ?, ?, ?)`

// CreateSession mints a fresh token for the user, valid for ttl.
func (d *DB) CreateSession(ctx context.Context, userID int64, userAgent string, ttl time.Duration) (Session, error) {
	s := Session{
		Token:     newSessionToken(),
		UserID:    userID,
		CreatedAt: nowMs(),
		UserAgent: userAgent,
	}
	s.ExpiresAt = s.CreatedAt + ttl.Milliseconds()
	_, err := d.write.ExecContext(ctx, createSessionSQL, s.Token, s.UserID, s.CreatedAt, s.ExpiresAt, s.UserAgent)
	return s, err
}

const sessionByTokenSQL = `SELECT s.token, s.user_id, s.created_at, s.expires_at, s.user_agent, ` + prefixedUserCols + `
	FROM sessions s JOIN users u ON u.id = s.user_id
	WHERE s.token = ? AND s.expires_at > ?`

const prefixedUserCols = "u.id, u.username, u.password_hash, u.display_name, u.avatar, u.is_admin, u.created_at, u.settings_json"

// SessionByToken resolves a live session and its user in one query. Expired
// or unknown tokens come back as ErrNotFound.
func (d *DB) SessionByToken(ctx context.Context, token string) (Session, User, error) {
	var s Session
	var u User
	err := d.read.QueryRowContext(ctx, sessionByTokenSQL, token, nowMs()).Scan(
		&s.Token, &s.UserID, &s.CreatedAt, &s.ExpiresAt, &s.UserAgent,
		&u.ID, &u.Username, &u.PasswordHash, &u.DisplayName, &u.Avatar, &u.IsAdmin, &u.CreatedAt, &u.SettingsJSON)
	return s, u, one(err)
}

const deleteSessionSQL = `DELETE FROM sessions WHERE token = ?`

// DeleteSession drops one session; deleting a missing token is not an error.
func (d *DB) DeleteSession(ctx context.Context, token string) error {
	_, err := d.write.ExecContext(ctx, deleteSessionSQL, token)
	return err
}

const deleteExpiredSessionsSQL = `DELETE FROM sessions WHERE expires_at <= ?`

// DeleteExpiredSessions removes every expired session and reports the count;
// the serve sweeper calls it hourly.
func (d *DB) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	res, err := d.write.ExecContext(ctx, deleteExpiredSessionsSQL, nowMs())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
