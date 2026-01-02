package definitions

import (
	"github.com/oriiyx/fritz/cmd/cli/config"
	"github.com/spf13/cobra"
)

func NewDefinitionsCmd(deps *config.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "definitions",
		Short: "Manage Fritz entity definitions",
		Long:  `Manage entity definitions for the Fritz system.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	// Add subcommands with dependencies
	cmd.AddCommand(NewLoadDefinitionsCmd(deps))

	return cmd
}
