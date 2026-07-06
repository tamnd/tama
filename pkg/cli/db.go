package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDBCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Inspect and migrate the database",
		Long: "Migrations also run automatically whenever tama opens the database;\n" +
			"these commands exist to run them explicitly and to see where the\n" +
			"schema stands.",
		Example: "  tama db status",
	}
	cmd.AddCommand(newDBMigrateCmd(), newDBStatusCmd())
	return cmd
}

func newDBMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate",
		Short:   "Apply pending migrations",
		Long:    "Applies every pending migration in order and prints each one it ran.",
		Example: "  tama db migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, db, err := openStore(cmd)
			if err != nil {
				return err
			}
			defer db.Close()

			applied := db.Applied()
			if len(applied) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "schema is up to date")
				return nil
			}
			for _, m := range applied {
				fmt.Fprintf(cmd.OutOrStdout(), "applied %04d_%s\n", m.Version, m.Name)
			}
			return nil
		},
	}
	addDataFlag(cmd)
	return cmd
}

func newDBStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Show the schema version",
		Long:    "Prints the current schema version and how many migrations are pending.",
		Example: "  tama db status",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, db, err := openStore(cmd)
			if err != nil {
				return err
			}
			defer db.Close()

			current, pending, err := db.SchemaVersion(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "schema version %d, %d pending\n", current, pending)
			return nil
		},
	}
	addDataFlag(cmd)
	return cmd
}
