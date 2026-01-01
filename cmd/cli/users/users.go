package users

import (
	"github.com/oriiyx/fritz/cmd/cli/config"
	"github.com/spf13/cobra"
)

func NewUsersCmd(deps *config.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage Fritz user accounts",
		Long: `Manage user accounts for the Fritz PIM system.

Provides commands for creating, listing, and managing user accounts.
All user operations require database access.`,
		Example: `  # Create a new user
  fritz users create user@example.com password123

  # List all users (future)
  fritz users list`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	// Add subcommands with dependencies
	cmd.AddCommand(NewCreateUserCmd(deps))

	return cmd
}
