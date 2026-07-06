// Package cli wires tama's command surface: the cobra tree, flags, and the
// fang-rendered help and errors. The real work lives in pkg/server, pkg/store,
// pkg/course, and friends; this layer parses flags and hands off.
package cli

import (
	"context"
	"fmt"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

// Version metadata, stamped by the release build through ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// Execute builds the root command and runs it through fang. main passes the
// signal-aware context so Ctrl-C shuts the server down cleanly. It returns the
// process exit code.
func Execute(ctx context.Context) int {
	root := newRoot()
	opts := []fang.Option{
		fang.WithVersion(Version),
	}
	if err := fang.Execute(ctx, root, opts...); err != nil {
		return 1
	}
	return 0
}

func newRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "tama",
		Short: "Self-hosted language learning with a path, hearts, streaks, and a cat",
		Long: "tama (タマ) is a language learning app in one binary: a lesson path with\n" +
			"hearts, streaks, XP, gems, leagues, and quests, for any pair of languages.\n" +
			"Course content is generated ahead of time by an LLM through an\n" +
			"OpenAI-compatible endpoint and frozen into course packs, so lessons are\n" +
			"instant and fully offline once a pack is cached.",
		Version:       fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date),
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newServeCmd())
	return root
}
