package cmd

import (
	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
	"github.com/user-cube/kontext/pkg/ui"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current Kubernetes context",
	Long: `Display the currently active Kubernetes context from your kubeconfig file.

Examples:
  # Show the current active context
  kontext current`,
	Run: func(cmd *cobra.Command, args []string) {
		currentContext, err := kubeconfig.GetCurrentContext()
		if err != nil {
			ui.PrintError("Error retrieving current context", err, true)
		}

		ui.PrintCurrentContext(currentContext)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
