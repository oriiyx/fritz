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
		Short: "Display Fritz CLI version information",
		Long: `Display version information for Fritz CLI.

Shows the version number, build date, git commit, Go version, and platform.
Use --json flag for machine-readable output or --short for version number only.`,
		Example: `  # Show full version details
  fritz version

  # Show only version number
  fritz version --short

  # Output as JSON
  fritz version --json`,
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
	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "output version information in JSON format")
	cmd.Flags().BoolVarP(&shortOutput, "short", "s", false, "output only the version number")

	return cmd
}
