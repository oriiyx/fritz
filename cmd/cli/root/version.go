package root

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oriiyx/fritz/cmd/cli/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print detailed version information about Fritz CLI including build date, git commit, and platform.`,
		Run: func(cmd *cobra.Command, args []string) {
			if jsonOutput {
				info := version.GetBuildInfo()
				output, err := json.MarshalIndent(info, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "%sError marshaling version info: %s%s\n", ColorRed, err.Error(), ColorReset)
					os.Exit(1)
				}
				fmt.Println(string(output))
			} else if shortOutput {
				fmt.Println(version.GetVersion())
			} else {
				fmt.Println(version.GetFullVersion())
			}
		},
	}

	// Version command flags
	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output version information in JSON format")
	cmd.Flags().BoolVarP(&shortOutput, "short", "s", false, "Output only the version number")

	return cmd
}
