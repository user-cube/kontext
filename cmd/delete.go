package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
	"github.com/user-cube/kontext/pkg/ui"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete [context]",
	Aliases: []string{"rm"},
	Short:   "Delete a Kubernetes context from your kubeconfig",
	Long: `Delete a Kubernetes context from your kubeconfig file.
If no context name is provided, an interactive selector will be displayed.

Examples:
  # Delete a context interactively
  kontext delete

  # Delete a specific context by name
  kontext delete my-context
  kontext rm my-context`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load available contexts
		contexts, err := kubeconfig.GetContexts()
		if err != nil {
			ui.PrintError("Error retrieving contexts", err, true)
		}

		if len(contexts) == 0 {
			ui.PrintWarning("No contexts found in kubeconfig")
			return
		}

		currentContext, err := kubeconfig.GetCurrentContext()
		if err != nil {
			ui.PrintError("Error retrieving current context", err, true)
		}

		var contextName string

		if len(args) == 0 {
			// Interactive selection if no context name is provided
			contextNames := make([]string, 0, len(contexts))
			for name := range contexts {
				contextNames = append(contextNames, name)
			}

			// Prioritize current context at the top
			contextNames = ui.SortContexts(contextNames, currentContext, true)

			selector := ui.CreateContextSelector(contextNames, currentContext)
			_, selection, err := selector.Run()
			if err != nil {
				ui.PrintError("Selection canceled", err, false)
				return
			}

			contextName = selection
		} else {
			contextName = args[0]

			// Validate that the context exists
			if _, exists := contexts[contextName]; !exists {
				contextNames := make([]string, 0, len(contexts))
				for name := range contexts {
					contextNames = append(contextNames, name)
				}
				contextNames = ui.SortContexts(contextNames, currentContext, false)

				ui.PrintError(fmt.Sprintf("Context '%s' does not exist", contextName), nil, false)
				ui.PrintContextList(contextNames, currentContext)
				os.Exit(1)
			}
		}

		// Warn if deleting the current context
		if contextName == currentContext {
			ui.PrintWarning(
				fmt.Sprintf("You are about to delete the current context '%s'", contextName),
				"(current context will be unset)",
			)
		}

		// Ask for confirmation before deleting
		confirmed, err := ui.ConfirmAction(fmt.Sprintf("Delete context '%s' from kubeconfig?", contextName))
		if err != nil {
			ui.PrintError("Error during confirmation", err, true)
		}
		if !confirmed {
			ui.PrintWarning("Context deletion canceled", contextName)
			return
		}

		// Perform deletion
		if err := kubeconfig.DeleteContext(contextName); err != nil {
			ui.PrintError("Error deleting context", err, true)
		}

		ui.PrintSuccess("Deleted context", contextName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

