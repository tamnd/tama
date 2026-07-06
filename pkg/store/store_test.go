package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func openTest(t *testing.T) *DB {
	t.Helper()
	db, err := Open(context.Background(), filepath.Join(t.TempDir(), "tama.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// freeze pins store.Now to a fixed instant and restores it afterwards.
func freeze(t *testing.T, at time.Time) {
	t.Helper()
	old := Now
	Now = func() time.Time { return at }
	t.Cleanup(func() { Now = old })
}

func mustCreateUser(t *testing.T, db *DB, name string) User {
	t.Helper()
	u, err := db.CreateUser(context.Background(), name, "$argon2id$fake", false)
	if err != nil {
		t.Fatalf("CreateUser(%s): %v", name, err)
	}
	return u
}

func mustCreateCourse(t *testing.T, db *DB, id string) {
	t.Helper()
	err := db.UpsertCourse(context.Background(), Course{
		ID: id, BaseLang: "en", TargetLang: "es", Title: "Spanish", Status: "ready",
	})
	if err != nil {
		t.Fatalf("UpsertCourse(%s): %v", id, err)
	}
}

func TestMigrateFromZero(t *testing.T) {
	db := openTest(t)
	if len(db.Applied()) == 0 {
		t.Fatal("fresh open applied no migrations")
	}
	current, pending, err := db.SchemaVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if current < 1 || pending != 0 {
		t.Errorf("schema version = %d with %d pending, want >=1 and 0", current, pending)
	}
}

func TestMigrateIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tama.db")
	ctx := context.Background()

	db, err := Open(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()

	db, err = Open(ctx, path)
	if err != nil {
		t.Fatalf("second Open: %v", err)
	}
	defer db.Close()
	if n := len(db.Applied()); n != 0 {
		t.Errorf("second open applied %d migrations, want 0", n)
	}
}

func TestCreateAndLookupUser(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()

	u := mustCreateUser(t, db, "mochi")
	if u.ID == 0 || u.CreatedAt == 0 || u.SettingsJSON != "{}" {
		t.Errorf("CreateUser row = %+v", u)
	}

	got, err := db.UserByUsername(ctx, "mochi")
	if err != nil || got.ID != u.ID {
		t.Errorf("UserByUsername = %+v, %v", got, err)
	}
	// The username column collates NOCASE.
	if _, err := db.UserByUsername(ctx, "MOCHI"); err != nil {
		t.Errorf("case-insensitive lookup: %v", err)
	}
	if _, err := db.UserByID(ctx, u.ID); err != nil {
		t.Errorf("UserByID: %v", err)
	}
	if _, err := db.UserByID(ctx, 999); !errors.Is(err, ErrNotFound) {
		t.Errorf("UserByID(999) = %v, want ErrNotFound", err)
	}

	if err := db.UpdatePassword(ctx, u.ID, "$argon2id$new"); err != nil {
		t.Errorf("UpdatePassword: %v", err)
	}
	got, _ = db.UserByID(ctx, u.ID)
	if got.PasswordHash != "$argon2id$new" {
		t.Errorf("hash not updated: %q", got.PasswordHash)
	}

	users, err := db.ListUsers(ctx)
	if err != nil || len(users) != 1 {
		t.Errorf("ListUsers = %v, %v", users, err)
	}
	if err := db.DeleteUser(ctx, u.ID); err != nil {
		t.Errorf("DeleteUser: %v", err)
	}
	if err := db.DeleteUser(ctx, u.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("second DeleteUser = %v, want ErrNotFound", err)
	}
}

func TestDuplicateUsernameFails(t *testing.T) {
	db := openTest(t)
	mustCreateUser(t, db, "mochi")
	if _, err := db.CreateUser(context.Background(), "Mochi", "$argon2id$fake", false); err == nil {
		t.Fatal("duplicate username (case-insensitive) inserted twice")
	}
}

func TestSessionLifecycle(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()
	u := mustCreateUser(t, db, "mochi")

	freeze(t, time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC))
	s, err := db.CreateSession(ctx, u.ID, "test-agent", 30*24*time.Hour)
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if len(s.Token) != 43 { // 32 bytes as unpadded base64url
		t.Errorf("token length = %d, want 43", len(s.Token))
	}

	got, gotUser, err := db.SessionByToken(ctx, s.Token)
	if err != nil || got.UserID != u.ID || gotUser.Username != "mochi" || got.UserAgent != "test-agent" {
		t.Errorf("SessionByToken = %+v, %+v, %v", got, gotUser, err)
	}

	// Move the clock past expiry: the token must stop resolving.
	freeze(t, time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC))
	if _, _, err := db.SessionByToken(ctx, s.Token); !errors.Is(err, ErrNotFound) {
		t.Errorf("expired session = %v, want ErrNotFound", err)
	}

	n, err := db.DeleteExpiredSessions(ctx)
	if err != nil || n != 1 {
		t.Errorf("DeleteExpiredSessions = %d, %v, want 1", n, err)
	}

	s2, _ := db.CreateSession(ctx, u.ID, "", time.Hour)
	if err := db.DeleteSession(ctx, s2.Token); err != nil {
		t.Errorf("DeleteSession: %v", err)
	}
	if _, _, err := db.SessionByToken(ctx, s2.Token); !errors.Is(err, ErrNotFound) {
		t.Errorf("deleted session still resolves: %v", err)
	}
}

func TestCoursesAndPacks(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()
	mustCreateCourse(t, db, "es-en")

	c, err := db.CourseByID(ctx, "es-en")
	if err != nil || c.Status != "ready" || c.PackID != nil {
		t.Errorf("CourseByID = %+v, %v", c, err)
	}

	// Upsert refreshes mutable fields without duplicating the row.
	if err := db.UpsertCourse(ctx, Course{ID: "es-en", BaseLang: "en", TargetLang: "es", Title: "Spanish 2", Status: "ready"}); err != nil {
		t.Fatalf("second UpsertCourse: %v", err)
	}
	all, err := db.ListCourses(ctx)
	if err != nil || len(all) != 1 || all[0].Title != "Spanish 2" {
		t.Errorf("ListCourses = %+v, %v", all, err)
	}

	if _, err := db.LatestPack(ctx, "es-en"); !errors.Is(err, ErrNotFound) {
		t.Errorf("LatestPack on packless course = %v, want ErrNotFound", err)
	}
	for v := 1; v <= 2; v++ {
		_, err := db.InsertPack(ctx, CoursePack{
			CourseID: "es-en", Version: v, Format: 1,
			Content: []byte{0x28, 0xb5, 0x2f, 0xfd, byte(v)}, SHA256: "deadbeef",
		})
		if err != nil {
			t.Fatalf("InsertPack v%d: %v", v, err)
		}
	}
	p, err := db.LatestPack(ctx, "es-en")
	if err != nil || p.Version != 2 {
		t.Errorf("LatestPack = %+v, %v, want version 2", p, err)
	}
	if _, err := db.InsertPack(ctx, CoursePack{CourseID: "es-en", Version: 2, Format: 1, Content: []byte{0}, SHA256: "x"}); err == nil {
		t.Error("duplicate (course, version) pack inserted twice")
	}
}

func TestUpsertProgress(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()
	u := mustCreateUser(t, db, "mochi")
	mustCreateCourse(t, db, "es-en")

	p := Progress{UserID: u.ID, CourseID: "es-en", NodeID: "u1-l1", Crowns: 1}
	if err := db.UpsertProgress(ctx, p); err != nil {
		t.Fatalf("UpsertProgress: %v", err)
	}
	done := int64(1_750_000_000_000)
	p.Crowns = 3
	p.CompletedAt = &done
	if err := db.UpsertProgress(ctx, p); err != nil {
		t.Fatalf("UpsertProgress update: %v", err)
	}

	rows, err := db.ProgressForCourse(ctx, u.ID, "es-en")
	if err != nil || len(rows) != 1 {
		t.Fatalf("ProgressForCourse = %+v, %v", rows, err)
	}
	if rows[0].Crowns != 3 || rows[0].CompletedAt == nil || *rows[0].CompletedAt != done {
		t.Errorf("progress row = %+v", rows[0])
	}
}

func TestXPLedger(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()
	u := mustCreateUser(t, db, "mochi")
	mustCreateCourse(t, db, "es-en")

	freeze(t, time.Date(2026, 7, 1, 8, 0, 0, 0, time.UTC))
	if err := db.InsertXPEvent(ctx, u.ID, "es-en", 10, "lesson"); err != nil {
		t.Fatalf("InsertXPEvent: %v", err)
	}
	freeze(t, time.Date(2026, 7, 2, 8, 0, 0, 0, time.UTC))
	if err := db.InsertXPEvent(ctx, u.ID, "", 5, "bonus"); err != nil {
		t.Fatalf("InsertXPEvent bonus: %v", err)
	}
	if err := db.InsertXPEvent(ctx, u.ID, "es-en", 1, "cheating"); err == nil {
		t.Error("xp reason outside the CHECK list accepted")
	}

	total, err := db.XPTotal(ctx, u.ID)
	if err != nil || total != 15 {
		t.Errorf("XPTotal = %d, %v, want 15", total, err)
	}
	since, err := db.XPSince(ctx, u.ID, time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC).UnixMilli())
	if err != nil || since != 5 {
		t.Errorf("XPSince = %d, %v, want 5", since, err)
	}
}

func TestStreaksHeartsGems(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()
	u := mustCreateUser(t, db, "mochi")

	if _, err := db.GetStreak(ctx, u.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("GetStreak before upsert = %v, want ErrNotFound", err)
	}
	if err := db.UpsertStreak(ctx, Streak{UserID: u.ID, Current: 3, Longest: 7, LastActiveDay: "2026-07-06"}); err != nil {
		t.Fatalf("UpsertStreak: %v", err)
	}
	s, err := db.GetStreak(ctx, u.ID)
	if err != nil || s.Current != 3 || s.Longest != 7 || s.LastActiveDay != "2026-07-06" {
		t.Errorf("GetStreak = %+v, %v", s, err)
	}

	refill := int64(1_750_000_000_000)
	if err := db.UpsertHearts(ctx, Hearts{UserID: u.ID, Count: 4, Max: 5, RefillStartedAt: &refill}); err != nil {
		t.Fatalf("UpsertHearts: %v", err)
	}
	h, err := db.GetHearts(ctx, u.ID)
	if err != nil || h.Count != 4 || h.RefillStartedAt == nil || *h.RefillStartedAt != refill || h.UnlimitedUntil != nil {
		t.Errorf("GetHearts = %+v, %v", h, err)
	}

	if err := db.InsertGemEvent(ctx, u.ID, 100, "quest", "q1"); err != nil {
		t.Fatalf("InsertGemEvent: %v", err)
	}
	if err := db.InsertGemEvent(ctx, u.ID, -30, "shop", "freeze"); err != nil {
		t.Fatalf("InsertGemEvent spend: %v", err)
	}
	bal, err := db.GemBalance(ctx, u.ID)
	if err != nil || bal != 70 {
		t.Errorf("GemBalance = %d, %v, want 70", bal, err)
	}
}

func TestDueReviewItems(t *testing.T) {
	db := openTest(t)
	ctx := context.Background()
	u := mustCreateUser(t, db, "mochi")
	mustCreateCourse(t, db, "es-en")

	now := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC).UnixMilli()
	for i, due := range []int64{now - 2000, now - 1000, now, now + 60_000} {
		item := ReviewItem{
			UserID: u.ID, CourseID: "es-en", ItemID: string(rune('a' + i)),
			Stability: 1.5, Difficulty: 5.0, Due: due,
		}
		if err := db.UpsertReviewItem(ctx, item); err != nil {
			t.Fatalf("UpsertReviewItem: %v", err)
		}
	}

	due, err := db.DueReviewItems(ctx, u.ID, "es-en", now, 10)
	if err != nil {
		t.Fatalf("DueReviewItems: %v", err)
	}
	if len(due) != 3 || due[0].ItemID != "a" || due[2].ItemID != "c" {
		t.Errorf("due queue = %+v, want a,b,c most overdue first", due)
	}

	limited, err := db.DueReviewItems(ctx, u.ID, "es-en", now, 2)
	if err != nil || len(limited) != 2 {
		t.Errorf("limited due queue = %+v, %v, want 2 items", limited, err)
	}

	// Rescheduling an item moves it out of the due window.
	if err := db.UpsertReviewItem(ctx, ReviewItem{UserID: u.ID, CourseID: "es-en", ItemID: "a", Due: now + 120_000, LastReview: &now, Lapses: 1}); err != nil {
		t.Fatalf("reschedule: %v", err)
	}
	due, _ = db.DueReviewItems(ctx, u.ID, "es-en", now, 10)
	if len(due) != 2 {
		t.Errorf("due after reschedule = %+v, want 2", due)
	}
}

func TestDemoPackFixture(t *testing.T) {
	if len(DemoPack()) == 0 {
		t.Fatal("demo pack fixture is empty")
	}
}
