package users

import (
	"github.com/oriiyx/fritz/cmd/cli/config"
	"github.com/spf13/cobra"
)

func NewUsersCmd(deps *config.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Handles user management",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	// Add subcommands with dependencies
	cmd.AddCommand(NewCreateUserCmd(deps))

	return cmd
}
