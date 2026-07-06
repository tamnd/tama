package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tamnd/tama/pkg/config"
	"github.com/tamnd/tama/pkg/server"
	"github.com/tamnd/tama/pkg/store"
)

func newServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the tama server",
		Long: "Starts the HTTP server with the embedded web UI. State goes to a single\n" +
			"SQLite file under the data directory, course packs next to it.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cmd.Flags())
			if err != nil {
				return err
			}

			st, err := store.Open(cfg.DBPath())
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}
			defer st.Close()

			srv := server.New(cfg, st)
			fmt.Printf("tama listening on http://localhost%s (data in %s)\n", cfg.Addr, cfg.DataDir)
			return srv.Run(cmd.Context())
		},
	}

	cmd.Flags().String("addr", ":4321", "listen address")
	cmd.Flags().String("data", "", "data directory (default ~/.tama)")
	return cmd
}
