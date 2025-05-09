package cmd

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
	Long:  `Display the version, build date, and git commit of your kontext installation.`,
	Run: func(cmd *cobra.Command, args []string) {
		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Printf("%s %s\n", bold("Kontext:"), cyan(Version))
		fmt.Printf("%s %s\n", bold("Git Commit:"), GitCommit)
		fmt.Printf("%s %s\n", bold("Built:"), BuildDate)
		fmt.Printf("%s %s/%s\n", bold("Platform:"), runtime.GOOS, runtime.GOARCH)
		fmt.Printf("%s %s\n", bold("Go Version:"), runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
