package store

import (
	"context"
	"database/sql"
)

// Streak mirrors the streaks table. LastActiveDay is the user-local
// YYYY-MM-DD, not a timestamp.
type Streak struct {
	UserID        int64
	Current       int
	Longest       int
	LastActiveDay string
	Freezes       int
	UpdatedAt     int64
}

// Hearts mirrors the hearts table; refill math lives in pkg/engine.
type Hearts struct {
	UserID          int64
	Count           int
	Max             int
	RefillStartedAt *int64
	UnlimitedUntil  *int64
}

const insertXPEventSQL = `INSERT INTO xp_events (user_id, course_id, amount, reason, created_at) VALUES (?, ?, ?, ?, ?)`

// InsertXPEvent appends to the XP ledger. courseID may be empty for
// course-agnostic awards and is stored as NULL.
func (d *DB) InsertXPEvent(ctx context.Context, userID int64, courseID string, amount int, reason string) error {
	var course any
	if courseID != "" {
		course = courseID
	}
	_, err := d.write.ExecContext(ctx, insertXPEventSQL, userID, course, amount, reason, nowMs())
	return err
}

const xpTotalSQL = `SELECT COALESCE(SUM(amount), 0) FROM xp_events WHERE user_id = ?`

// XPTotal sums the whole ledger for a user.
func (d *DB) XPTotal(ctx context.Context, userID int64) (int64, error) {
	var total int64
	err := d.read.QueryRowContext(ctx, xpTotalSQL, userID).Scan(&total)
	return total, err
}

const xpSinceSQL = `SELECT COALESCE(SUM(amount), 0) FROM xp_events WHERE user_id = ? AND created_at >= ?`

// XPSince sums XP earned at or after the given unix-millisecond instant.
func (d *DB) XPSince(ctx context.Context, userID int64, sinceMs int64) (int64, error) {
	var total int64
	err := d.read.QueryRowContext(ctx, xpSinceSQL, userID, sinceMs).Scan(&total)
	return total, err
}

const getStreakSQL = `SELECT user_id, current, longest, last_active_day, freezes, updated_at FROM streaks WHERE user_id = ?`

// GetStreak fetches the streak row, ErrNotFound before the first upsert.
func (d *DB) GetStreak(ctx context.Context, userID int64) (Streak, error) {
	var s Streak
	err := d.read.QueryRowContext(ctx, getStreakSQL, userID).Scan(
		&s.UserID, &s.Current, &s.Longest, &s.LastActiveDay, &s.Freezes, &s.UpdatedAt)
	return s, one(err)
}

const upsertStreakSQL = `INSERT INTO streaks (user_id, current, longest, last_active_day, freezes, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT (user_id) DO UPDATE SET
		current = excluded.current,
		longest = excluded.longest,
		last_active_day = excluded.last_active_day,
		freezes = excluded.freezes,
		updated_at = excluded.updated_at`

// UpsertStreak writes the whole streak row.
func (d *DB) UpsertStreak(ctx context.Context, s Streak) error {
	_, err := d.write.ExecContext(ctx, upsertStreakSQL,
		s.UserID, s.Current, s.Longest, s.LastActiveDay, s.Freezes, nowMs())
	return err
}

const getHeartsSQL = `SELECT user_id, count, max, refill_started_at, unlimited_until FROM hearts WHERE user_id = ?`

// GetHearts fetches the hearts row, ErrNotFound before the first upsert.
func (d *DB) GetHearts(ctx context.Context, userID int64) (Hearts, error) {
	var h Hearts
	var refill, unlimited sql.NullInt64
	err := d.read.QueryRowContext(ctx, getHeartsSQL, userID).Scan(&h.UserID, &h.Count, &h.Max, &refill, &unlimited)
	if refill.Valid {
		h.RefillStartedAt = &refill.Int64
	}
	if unlimited.Valid {
		h.UnlimitedUntil = &unlimited.Int64
	}
	return h, one(err)
}

const upsertHeartsSQL = `INSERT INTO hearts (user_id, count, max, refill_started_at, unlimited_until)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT (user_id) DO UPDATE SET
		count = excluded.count,
		max = excluded.max,
		refill_started_at = excluded.refill_started_at,
		unlimited_until = excluded.unlimited_until`

// UpsertHearts writes the whole hearts row.
func (d *DB) UpsertHearts(ctx context.Context, h Hearts) error {
	_, err := d.write.ExecContext(ctx, upsertHeartsSQL, h.UserID, h.Count, h.Max, h.RefillStartedAt, h.UnlimitedUntil)
	return err
}

const insertGemEventSQL = `INSERT INTO gems_ledger (user_id, amount, reason, ref, created_at) VALUES (?, ?, ?, ?, ?)`

// InsertGemEvent appends a signed amount to the gems ledger. Overdraft checks
// belong to the engine, inside its write transaction.
func (d *DB) InsertGemEvent(ctx context.Context, userID int64, amount int, reason, ref string) error {
	_, err := d.write.ExecContext(ctx, insertGemEventSQL, userID, amount, reason, ref, nowMs())
	return err
}

const gemBalanceSQL = `SELECT COALESCE(SUM(amount), 0) FROM gems_ledger WHERE user_id = ?`

// GemBalance sums the gems ledger for a user.
func (d *DB) GemBalance(ctx context.Context, userID int64) (int64, error) {
	var total int64
	err := d.read.QueryRowContext(ctx, gemBalanceSQL, userID).Scan(&total)
	return total, err
}
