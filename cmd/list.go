package cmd

import (
	"sort"

	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
	"github.com/user-cube/kontext/pkg/ui"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Kubernetes contexts",
	Long: `List all available Kubernetes contexts from your kubeconfig file.

Examples:
  # List all available contexts with the current one highlighted
  kontext list`,
	Run: func(cmd *cobra.Command, args []string) {
		contexts, err := kubeconfig.GetContexts()
		if err != nil {
			ui.PrintError("Error retrieving contexts", err, true)
		}

		currentContext, err := kubeconfig.GetCurrentContext()
		if err != nil {
			ui.PrintError("Error retrieving current context", err, true)
		}

		// Sort context names for consistent display
		contextNames := make([]string, 0, len(contexts))
		for name := range contexts {
			contextNames = append(contextNames, name)
		}
		sort.Strings(contextNames)

		// Print contexts using the UI package
		ui.PrintContextList(contextNames, currentContext)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
