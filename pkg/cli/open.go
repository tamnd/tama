package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tamnd/tama/pkg/config"
	"github.com/tamnd/tama/pkg/store"
)

// addDataFlag registers --data on commands that open the database without
// running the server.
func addDataFlag(cmd *cobra.Command) {
	cmd.Flags().String("data", "", "data directory (default ~/.tama)")
}

// openStore resolves config and opens the database, applying any pending
// migrations on the way, same as serve does.
func openStore(cmd *cobra.Command) (*config.Config, *store.DB, error) {
	cfg, err := config.Load(cmd.Flags())
	if err != nil {
		return nil, nil, err
	}
	db, err := store.Open(cmd.Context(), cfg.DBPath())
	if err != nil {
		return nil, nil, fmt.Errorf("open store: %w", err)
	}
	return cfg, db, nil
}
