package store

import (
	"context"
	"database/sql"
)

// ReviewItem mirrors the review_items table, the FSRS state per item. The
// scheduler that fills it lands with pkg/engine in M6.
type ReviewItem struct {
	UserID     int64
	CourseID   string
	ItemID     string
	State      int
	Stability  float64
	Difficulty float64
	Due        int64
	LastReview *int64
	Lapses     int
}

const upsertReviewItemSQL = `INSERT INTO review_items (user_id, course_id, item_id, state, stability, difficulty, due, last_review, lapses)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (user_id, course_id, item_id) DO UPDATE SET
		state = excluded.state,
		stability = excluded.stability,
		difficulty = excluded.difficulty,
		due = excluded.due,
		last_review = excluded.last_review,
		lapses = excluded.lapses`

// UpsertReviewItem writes the full scheduling state of one item.
func (d *DB) UpsertReviewItem(ctx context.Context, r ReviewItem) error {
	_, err := d.write.ExecContext(ctx, upsertReviewItemSQL,
		r.UserID, r.CourseID, r.ItemID, r.State, r.Stability, r.Difficulty, r.Due, r.LastReview, r.Lapses)
	return err
}

const dueReviewItemsSQL = `SELECT user_id, course_id, item_id, state, stability, difficulty, due, last_review, lapses
	FROM review_items WHERE user_id = ? AND course_id = ? AND due <= ?
	ORDER BY due LIMIT ?`

// DueReviewItems returns up to limit items due at or before now, most
// overdue first.
func (d *DB) DueReviewItems(ctx context.Context, userID int64, courseID string, now int64, limit int) ([]ReviewItem, error) {
	rows, err := d.read.QueryContext(ctx, dueReviewItemsSQL, userID, courseID, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ReviewItem
	for rows.Next() {
		var r ReviewItem
		var last sql.NullInt64
		if err := rows.Scan(&r.UserID, &r.CourseID, &r.ItemID, &r.State, &r.Stability, &r.Difficulty, &r.Due, &last, &r.Lapses); err != nil {
			return nil, err
		}
		if last.Valid {
			r.LastReview = &last.Int64
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
