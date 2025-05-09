package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
	"github.com/user-cube/kontext/pkg/ui"
)

// nsCmd represents the namespace command
var nsCmd = &cobra.Command{
	Use:     "namespace [namespace]",
	Aliases: []string{"ns"},
	Short:   "View or change the current namespace",
	Long: `View or change the namespace for the current Kubernetes context.
If no namespace is provided, an interactive selection menu will be displayed.

Examples:
  # Show interactive namespace selector 
  kontext namespace
  
  # Use short alias to show interactive namespace selector
  kontext ns
  
  # Show current namespace without interactive selector
  kontext namespace --show
  kontext ns -s
  
  # Switch to a specific namespace directly
  kontext namespace my-namespace
  kontext ns my-namespace
  
  # Typical workflow: switch context, then namespace
  kontext switch my-context
  kontext ns my-namespace`,
	Run: runNamespace,
}

// Function to run the namespace command
func runNamespace(cmd *cobra.Command, args []string) {
	// Get current context and namespace
	currentContext, err := kubeconfig.GetCurrentContext()
	if err != nil {
		ui.PrintError("Error retrieving current context", err, true)
	}

	currentNamespace, err := kubeconfig.GetCurrentNamespace()
	if err != nil {
		ui.PrintError("Error retrieving current namespace", err, true)
	}

	// If no arguments are provided, show the current namespace or interactive selector
	if len(args) == 0 {
		// Check if a specific flag was provided to only show the current value
		if showOnly, _ := cmd.Flags().GetBool("show"); showOnly {
			ui.PrintCurrentNamespace(currentContext, currentNamespace)
			return
		}

		// Get available namespaces
		namespaces, err := kubeconfig.GetNamespaces()
		if err != nil {
			ui.PrintError("Error retrieving namespaces", err, true)
		}

		// Check if we have namespaces to display
		if len(namespaces) == 0 {
			ui.PrintWarning("No namespaces available for context", currentContext)
			return
		}

		// Sort namespaces and prioritize the current namespace
		namespaces = ui.SortNamespaces(namespaces, currentNamespace, true)

		// Create an interactive selector
		selector := ui.CreateNamespaceSelector(namespaces, currentNamespace, currentContext)
		_, selection, err := selector.Run()

		if err != nil {
			ui.PrintError("Selection canceled", err, false)
			return
		}

		// If the selected namespace is the same as the current one, don't do anything
		if selection == currentNamespace {
			ui.PrintWarning(fmt.Sprintf("Namespace '%s' is already selected", selection))
			return
		}

		// Change the namespace
		err = kubeconfig.SetNamespace(selection)
		if err != nil {
			ui.PrintError("Error setting namespace", err, true)
		}

		ui.PrintSuccess("Switched to namespace", selection, fmt.Sprintf("in context %s", currentContext))
		return
	}

	// Change to the specified namespace
	namespace := args[0]

	// If the specified namespace is the same as the current one, don't do anything
	if namespace == currentNamespace {
		ui.PrintWarning(fmt.Sprintf("Namespace '%s' is already selected", namespace))
		return
	}

	// Verify that the specified namespace exists for this context
	namespaces, err := kubeconfig.GetNamespaces()
	if err != nil {
		ui.PrintError("Error retrieving namespaces", err, true)
	}

	// Check if the specified namespace exists
	namespaceExists := false
	for _, ns := range namespaces {
		if ns == namespace {
			namespaceExists = true
			break
		}
	}

	if !namespaceExists {
		ui.PrintWarning(fmt.Sprintf("Namespace '%s' does not exist in context '%s'", namespace, currentContext))
		// Continue anyway since the user explicitly requested this namespace
	}

	// Change the namespace
	err = kubeconfig.SetNamespace(namespace)
	if err != nil {
		ui.PrintError("Error setting namespace", err, true)
	}

	ui.PrintSuccess("Switched to namespace", namespace, fmt.Sprintf("in context %s", currentContext))
}

func init() {
	rootCmd.AddCommand(nsCmd)

	// Add flags
	nsCmd.Flags().BoolP("show", "s", false, "Only show the current namespace without the selector")
}
