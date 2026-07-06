package cli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/klauspost/compress/zstd"
	"github.com/spf13/cobra"

	"github.com/tamnd/tama/pkg/api"
	"github.com/tamnd/tama/pkg/store"
)

func newSeedCmd() *cobra.Command {
	var demo bool
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed demo data",
		Long: "Seeds the demo user (demo / demo1234), a ready es-from-en course with a tiny\n" +
			"handcrafted pack, a progress row, an xp event, and a 3-day streak, so\n" +
			"the app is explorable right after boot. Safe to run twice.",
		Example: "  tama seed --demo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !demo {
				return fmt.Errorf("nothing to seed, pass --demo")
			}
			_, db, err := openStore(cmd)
			if err != nil {
				return err
			}
			defer db.Close()

			if err := seedDemo(cmd.Context(), db); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "seeded demo user and es-from-en course (login demo / demo1234)")
			return nil
		},
	}
	cmd.Flags().BoolVar(&demo, "demo", false, "seed the demo dataset")
	addDataFlag(cmd)
	return cmd
}

// seedDemo is idempotent: every step either upserts or checks before it
// writes, so a second run changes nothing.
func seedDemo(ctx context.Context, db *store.DB) error {
	user, err := db.UserByUsername(ctx, "demo")
	if errors.Is(err, store.ErrNotFound) {
		hash, err := api.HashPassword("demo1234")
		if err != nil {
			return err
		}
		if user, err = db.CreateUser(ctx, "demo", hash, false); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err := db.UpsertCourse(ctx, store.Course{
		ID: "es-from-en", BaseLang: "en", TargetLang: "es",
		Title: "Spanish from English", Status: "ready",
	}); err != nil {
		return err
	}

	packID, err := seedDemoPack(ctx, db)
	if err != nil {
		return err
	}
	if err := db.UpsertCourse(ctx, store.Course{
		ID: "es-from-en", BaseLang: "en", TargetLang: "es",
		Title: "Spanish from English", Status: "ready", PackID: &packID,
	}); err != nil {
		return err
	}

	if err := db.UpsertProgress(ctx, store.Progress{
		UserID: user.ID, CourseID: "es-from-en", NodeID: "u1-l1", Crowns: 1,
	}); err != nil {
		return err
	}

	total, err := db.XPTotal(ctx, user.ID)
	if err != nil {
		return err
	}
	if total == 0 {
		if err := db.InsertXPEvent(ctx, user.ID, "es-from-en", 10, "lesson"); err != nil {
			return err
		}
	}

	return db.UpsertStreak(ctx, store.Streak{
		UserID: user.ID, Current: 3, Longest: 3,
		LastActiveDay: store.Now().Format("2006-01-02"),
	})
}

// seedDemoPack compresses the checked-in fixture into course_packs, reusing
// the existing row on re-runs.
func seedDemoPack(ctx context.Context, db *store.DB) (int64, error) {
	if p, err := db.LatestPack(ctx, "es-from-en"); err == nil {
		return p.ID, nil
	} else if !errors.Is(err, store.ErrNotFound) {
		return 0, err
	}

	raw := store.DemoPack()
	enc, err := zstd.NewWriter(nil)
	if err != nil {
		return 0, err
	}
	compressed := enc.EncodeAll(raw, nil)
	enc.Close()
	sum := sha256.Sum256(raw)

	return db.InsertPack(ctx, store.CoursePack{
		CourseID: "es-from-en", Version: 1, Format: 1,
		Content: compressed, SHA256: hex.EncodeToString(sum[:]),
		GeneratedBy: "fixture",
	})
}
