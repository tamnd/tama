package store

import (
	"context"
	"database/sql"
)

// Course mirrors the courses table. ID is "<target>-<base>" like "es-en".
type Course struct {
	ID         string
	BaseLang   string
	TargetLang string
	Title      string
	PackID     *int64
	Status     string
	CreatedAt  int64
}

// CoursePack mirrors the course_packs table. Content is the zstd-compressed
// pack JSON.
type CoursePack struct {
	ID          int64
	CourseID    string
	Version     int
	Format      int
	Content     []byte
	SHA256      string
	GeneratedBy string
	CreatedAt   int64
}

const upsertCourseSQL = `INSERT INTO courses (id, base_lang, target_lang, title, pack_id, status, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (id) DO UPDATE SET title = excluded.title, pack_id = excluded.pack_id, status = excluded.status`

// UpsertCourse inserts the course row or refreshes its title, pack pointer,
// and status; id and language pair are immutable once created.
func (d *DB) UpsertCourse(ctx context.Context, c Course) error {
	_, err := d.write.ExecContext(ctx, upsertCourseSQL, c.ID, c.BaseLang, c.TargetLang, c.Title, c.PackID, c.Status, nowMs())
	return err
}

const courseCols = "id, base_lang, target_lang, title, pack_id, status, created_at"

func scanCourse(row interface{ Scan(...any) error }) (Course, error) {
	var c Course
	var packID sql.NullInt64
	err := row.Scan(&c.ID, &c.BaseLang, &c.TargetLang, &c.Title, &packID, &c.Status, &c.CreatedAt)
	if packID.Valid {
		c.PackID = &packID.Int64
	}
	return c, one(err)
}

const courseByIDSQL = `SELECT ` + courseCols + ` FROM courses WHERE id = ?`

// CourseByID fetches one course row, ErrNotFound when absent.
func (d *DB) CourseByID(ctx context.Context, id string) (Course, error) {
	return scanCourse(d.read.QueryRowContext(ctx, courseByIDSQL, id))
}

const listCoursesSQL = `SELECT ` + courseCols + ` FROM courses ORDER BY id`

// ListCourses returns every course row.
func (d *DB) ListCourses(ctx context.Context) ([]Course, error) {
	rows, err := d.read.QueryContext(ctx, listCoursesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		c, err := scanCourse(rows)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

const insertPackSQL = `INSERT INTO course_packs (course_id, version, format, content, sha256, generated_by, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`

// InsertPack stores one immutable pack version and returns its row id.
func (d *DB) InsertPack(ctx context.Context, p CoursePack) (int64, error) {
	var id int64
	err := d.write.QueryRowContext(ctx, insertPackSQL,
		p.CourseID, p.Version, p.Format, p.Content, p.SHA256, p.GeneratedBy, nowMs()).Scan(&id)
	return id, err
}

const latestPackSQL = `SELECT id, course_id, version, format, content, sha256, generated_by, created_at
	FROM course_packs WHERE course_id = ? ORDER BY version DESC LIMIT 1`

// LatestPack returns the newest pack for a course, ErrNotFound when the
// course has none yet.
func (d *DB) LatestPack(ctx context.Context, courseID string) (CoursePack, error) {
	var p CoursePack
	err := d.read.QueryRowContext(ctx, latestPackSQL, courseID).Scan(
		&p.ID, &p.CourseID, &p.Version, &p.Format, &p.Content, &p.SHA256, &p.GeneratedBy, &p.CreatedAt)
	return p, one(err)
}
