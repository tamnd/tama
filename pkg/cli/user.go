package cli

import (
	"errors"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/tamnd/tama/pkg/api"
	"github.com/tamnd/tama/pkg/store"
)

func newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage user accounts",
		Long: "Manages the users table directly, no server required. Handy for the\n" +
			"first admin account and for resetting a forgotten password.",
		Example: "  tama user add --username mochi --password 'correct horse'",
	}
	cmd.AddCommand(newUserAddCmd(), newUserLsCmd(), newUserRmCmd(), newUserPasswdCmd())
	return cmd
}

func newUserAddCmd() *cobra.Command {
	var (
		username string
		password string
		admin    bool
	)
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Create a user",
		Long: "Creates a user with an argon2id-hashed password. Usernames are 3-24\n" +
			"characters of a-z, 0-9, or _ and are stored lowercased.",
		Example: "  tama user add --username mochi --password 'correct horse' --admin",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := api.ValidateUsername(username)
			if err != nil {
				return err
			}
			if err := api.ValidatePassword(password); err != nil {
				return err
			}
			hash, err := api.HashPassword(password)
			if err != nil {
				return err
			}

			_, db, err := openStore(cmd)
			if err != nil {
				return err
			}
			defer db.Close()

			u, err := db.CreateUser(cmd.Context(), name, hash, admin)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "created user %s (id %d)\n", u.Username, u.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "username (required)")
	cmd.Flags().StringVar(&password, "password", "", "password (required)")
	cmd.Flags().BoolVar(&admin, "admin", false, "grant admin")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	addDataFlag(cmd)
	return cmd
}

func newUserLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List users",
		Long:    "Lists every user with id, admin bit, and creation time.",
		Example: "  tama user ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, db, err := openStore(cmd)
			if err != nil {
				return err
			}
			defer db.Close()

			users, err := db.ListUsers(cmd.Context())
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tUSERNAME\tADMIN\tCREATED")
			for _, u := range users {
				fmt.Fprintf(w, "%d\t%s\t%v\t%s\n", u.ID, u.Username, u.IsAdmin,
					time.UnixMilli(u.CreatedAt).UTC().Format("2006-01-02 15:04"))
			}
			return w.Flush()
		},
	}
	addDataFlag(cmd)
	return cmd
}

func newUserRmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm <username>",
		Short:   "Delete a user",
		Long:    "Deletes a user; their sessions and progress cascade away with them.",
		Example: "  tama user rm mochi",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, db, err := openStore(cmd)
			if err != nil {
				return err
			}
			defer db.Close()

			u, err := db.UserByUsername(cmd.Context(), args[0])
			if errors.Is(err, store.ErrNotFound) {
				return fmt.Errorf("no such user %q", args[0])
			}
			if err != nil {
				return err
			}
			if err := db.DeleteUser(cmd.Context(), u.ID); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted user %s\n", u.Username)
			return nil
		},
	}
	addDataFlag(cmd)
	return cmd
}

func newUserPasswdCmd() *cobra.Command {
	var password string
	cmd := &cobra.Command{
		Use:     "passwd <username>",
		Short:   "Change a user's password",
		Long:    "Replaces the stored hash with an argon2id hash of the new password.",
		Example: "  tama user passwd mochi --password 'new horse'",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := api.ValidatePassword(password); err != nil {
				return err
			}
			hash, err := api.HashPassword(password)
			if err != nil {
				return err
			}

			_, db, err := openStore(cmd)
			if err != nil {
				return err
			}
			defer db.Close()

			u, err := db.UserByUsername(cmd.Context(), args[0])
			if errors.Is(err, store.ErrNotFound) {
				return fmt.Errorf("no such user %q", args[0])
			}
			if err != nil {
				return err
			}
			if err := db.UpdatePassword(cmd.Context(), u.ID, hash); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "updated password for %s\n", u.Username)
			return nil
		},
	}
	cmd.Flags().StringVar(&password, "password", "", "new password (required)")
	cmd.MarkFlagRequired("password")
	addDataFlag(cmd)
	return cmd
}
