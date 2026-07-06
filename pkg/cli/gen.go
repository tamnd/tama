package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tamnd/tama/pkg/config"
	"github.com/tamnd/tama/pkg/gen"
)

func newGenCmd() *cobra.Command {
	var (
		course string
		dryRun bool
	)
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate course packs",
		Long: "Course pack generation through the configured OpenAI-compatible\n" +
			"endpoint. M1 only proves connectivity: --dry-run pings the endpoint\n" +
			"and lists its models. Real generation lands in M8.",
		Example: "  tama gen --course es-en --dry-run",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !dryRun {
				return fmt.Errorf("generation is not implemented yet, run with --dry-run to test connectivity")
			}
			cfg, err := config.Load(cmd.Flags())
			if err != nil {
				return err
			}

			client := gen.New(cfg.LLM)
			models, err := client.Ping(cmd.Context())
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "endpoint %s is reachable, %d models:\n", cfg.LLM.BaseURL, len(models))
			for _, m := range models {
				fmt.Fprintf(out, "  %s\n", m.ID)
			}
			if course != "" {
				fmt.Fprintf(out, "would generate a pack for %s\n", course)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&course, "course", "", "course pair like es-en")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "only check endpoint connectivity")
	addDataFlag(cmd)
	return cmd
}
