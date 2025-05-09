package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/ui"
)

// Version information
var (
	Version   = "v0.1.0"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information of kontext",
	Long: `Display the version, build date, and git commit of your kontext installation.

Examples:
  # Show version information
  kontext version`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintInfo("Kontext", Version)
		ui.PrintInfo("Git Commit", GitCommit)
		ui.PrintInfo("Built", BuildDate)
		ui.PrintInfo("Platform", runtime.GOOS+"/"+runtime.GOARCH)
		ui.PrintInfo("Go Version", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
