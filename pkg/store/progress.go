package store

import (
	"context"
	"database/sql"
)

// Progress mirrors the progress table, one row per path node a user touched.
type Progress struct {
	UserID      int64
	CourseID    string
	NodeID      string
	Crowns      int
	Legendary   bool
	CompletedAt *int64
	UpdatedAt   int64
}

const upsertProgressSQL = `INSERT INTO progress (user_id, course_id, node_id, crowns, legendary, completed_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (user_id, course_id, node_id) DO UPDATE SET
		crowns = excluded.crowns,
		legendary = excluded.legendary,
		completed_at = excluded.completed_at,
		updated_at = excluded.updated_at`

// UpsertProgress writes the current state of one node for one user.
func (d *DB) UpsertProgress(ctx context.Context, p Progress) error {
	_, err := d.write.ExecContext(ctx, upsertProgressSQL,
		p.UserID, p.CourseID, p.NodeID, p.Crowns, p.Legendary, p.CompletedAt, nowMs())
	return err
}

const progressForCourseSQL = `SELECT user_id, course_id, node_id, crowns, legendary, completed_at, updated_at
	FROM progress WHERE user_id = ? AND course_id = ? ORDER BY node_id`

// ProgressForCourse returns every node row the user has for the course.
func (d *DB) ProgressForCourse(ctx context.Context, userID int64, courseID string) ([]Progress, error) {
	rows, err := d.read.QueryContext(ctx, progressForCourseSQL, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Progress
	for rows.Next() {
		var p Progress
		var completed sql.NullInt64
		if err := rows.Scan(&p.UserID, &p.CourseID, &p.NodeID, &p.Crowns, &p.Legendary, &completed, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if completed.Valid {
			p.CompletedAt = &completed.Int64
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
