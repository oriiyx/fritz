package definitions

import (
	"github.com/oriiyx/fritz/cmd/cli/config"
	"github.com/spf13/cobra"
)

func NewLoadDefinitionsCmd(deps *config.Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "load",
		Short: "Loads definitions into Fritz system",
		Long:  `Load definitions from .json files located inside var/entities/definitions`,
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// todo - handle storing current definition into database so that we have a reference point for comparing?
			return nil
		},
	}
}
