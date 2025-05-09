package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
	"github.com/user-cube/kontext/pkg/ui"
)

// runSwitch contains the main logic for the switch command
// This is shared with the root command to enable the same behavior with `kontext`
func runSwitch(cmd *cobra.Command, args []string) {
	var contextName string

	// Get available contexts
	contexts, err := kubeconfig.GetContexts()
	if err != nil {
		ui.PrintError("Error retrieving contexts", err, true)
	}

	currentContext, err := kubeconfig.GetCurrentContext()
	if err != nil {
		ui.PrintError("Error retrieving current context", err, true)
	}

	// Get the current namespace
	currentNamespace, err := kubeconfig.GetCurrentNamespace()
	if err != nil {
		// Non-fatal - default to "default" namespace
		currentNamespace = "default"
	}

	if len(args) == 0 {
		// If no context is provided, show interactive selector
		contextNames := make([]string, 0, len(contexts))
		for name := range contexts {
			contextNames = append(contextNames, name)
		}

		// Sort context names and optionally prioritize current context
		// Setting the third parameter to true would place current context first
		// Setting it to false maintains alphabetical order
		contextNames = ui.SortContexts(contextNames, currentContext, true)

		// Create the selector and run it
		selector := ui.CreateContextSelector(contextNames, currentContext)
		_, selection, err := selector.Run()

		if err != nil {
			ui.PrintError("Selection canceled", err, false)
			return
		}

		contextName = selection

		// Don't switch if selected context is already current
		if contextName == currentContext {
			ui.PrintWarning(fmt.Sprintf("Context '%s' is already selected", contextName))
			ui.PrintCurrentNamespace(contextName, currentNamespace)
			return
		}

		// Switch to the selected context
		err = kubeconfig.SwitchContext(contextName)
		if err != nil {
			ui.PrintError("Error switching context", err, true)
		}

		ui.PrintSuccess("Switched to context", contextName)

		// Get namespace for the new context
		targetNamespace, err := kubeconfig.GetNamespaceForContext(contextName)
		if err != nil {
			// Non-fatal - default to "default" namespace
			targetNamespace = "default"
		}

		ui.PrintSuccess("Namespace", targetNamespace)

		// The namespace selector will be handled by the caller if needed
		// We don't want to call it here to avoid duplicate namespace selection
	} else {
		contextName = args[0]

		// Check if the context exists
		if _, exists := contexts[contextName]; !exists {
			// Get available context names for the error message
			contextNames := make([]string, 0, len(contexts))
			for name := range contexts {
				contextNames = append(contextNames, name)
			}
			// Sort contexts with the current context highlighted
			contextNames = ui.SortContexts(contextNames, currentContext, false)

			ui.PrintError(fmt.Sprintf("Context '%s' does not exist", contextName), nil, false)
			ui.PrintContextList(contextNames, currentContext)
			os.Exit(1)
		}

		// Don't switch if selected context is already current
		if contextName == currentContext {
			ui.PrintWarning(fmt.Sprintf("Context '%s' is already selected", contextName))
			ui.PrintCurrentNamespace(contextName, currentNamespace)
			return
		}

		// Switch to the selected context
		err = kubeconfig.SwitchContext(contextName)
		if err != nil {
			ui.PrintError("Error switching context", err, true)
		}

		ui.PrintSuccess("Switched to context", contextName)

		// Get namespace for the new context
		targetNamespace, err := kubeconfig.GetNamespaceForContext(contextName)
		if err != nil {
			// Non-fatal - default to "default" namespace
			targetNamespace = "default"
		}

		ui.PrintSuccess("Namespace", targetNamespace)

		// The namespace selector will be handled by the caller if needed
		// We don't want to call it here to avoid duplicate namespace selection
	}
}

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [context] [namespace]",
	Short: "Switch to a specific Kubernetes context",
	Long: `Switch to a specific Kubernetes context in your kubeconfig file.
If no context is provided, an interactive selection menu will be displayed.

Examples:
  # Show interactive context selector
  kontext switch
  
  # Switch to specific context by name
  kontext switch my-context
  
  # Switch to context and then select namespace interactively
  kontext switch -n
  
  # Switch to specific context and then select namespace interactively
  kontext switch my-context -n
  
  # Switch to specific context and directly set a namespace
  kontext switch my-context -n my-namespace
  
  # The root command also acts as an alias to switch
  kontext
  kontext my-context
  kontext -n
  kontext my-context -n
  kontext my-context -n my-namespace`,
	ValidArgsFunction: contextCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the list of non-flag arguments (context and possibly namespace)
		nonFlagArgs := []string{}
		for _, arg := range args {
			// Check if this arg looks like a flag
			if !strings.HasPrefix(arg, "-") {
				nonFlagArgs = append(nonFlagArgs, arg)
			}
		}

		// Check if a namespace is specified after the -n flag
		setNS, _ := cmd.Flags().GetBool("set-namespace")
		var namespaceArg string

		// Get all args that might contain the -n flag and a namespace argument
		allArgs := os.Args
		if setNS {
			// Look for a namespace argument after the -n flag
			for i, arg := range allArgs {
				if (arg == "-n" || arg == "--set-namespace") && i+1 < len(allArgs) && !strings.HasPrefix(allArgs[i+1], "-") {
					// Check if the next arg isn't another flag and isn't part of the context name
					if len(nonFlagArgs) == 0 || allArgs[i+1] != nonFlagArgs[0] {
						namespaceArg = allArgs[i+1]
						break
					}
				}
			}
		}

		// If we found a namespace arg, run with it; otherwise use the normal runSwitch
		if namespaceArg != "" {
			// If we have a context, switch to it first, then set the namespace directly
			if len(nonFlagArgs) > 0 {
				// Create args with just the context
				contextArgs := []string{nonFlagArgs[0]}
				// Call runSwitch with just the context
				runSwitch(cmd, contextArgs)
				
				// Before setting the namespace, check if it exists
				newContext, err := kubeconfig.GetCurrentContext()
				if err != nil {
					ui.PrintError("Error retrieving current context", err, true)
				}
				
				namespaces, err := kubeconfig.GetNamespaces()
				if err != nil {
					ui.PrintError("Error retrieving namespaces", err, true)
				}
				
				// Check if the specified namespace exists
				namespaceExists := false
				for _, ns := range namespaces {
					if ns == namespaceArg {
						namespaceExists = true
						break
					}
				}
				
				if !namespaceExists {
					ui.PrintWarning(fmt.Sprintf("Namespace '%s' does not exist in context '%s'", namespaceArg, newContext))
					// Continue anyway since the user explicitly requested this namespace
				}
				
				// Then set the namespace directly
				nsCmd.Run(nsCmd, []string{namespaceArg})
			} else {
				// If no context specified, show context selector and then set namespace
				runSwitch(cmd, []string{})
				nsCmd.Run(nsCmd, []string{namespaceArg})
			}
		} else if setNS {
			// No namespace argument but -n flag was specified
			// First switch context if needed
			runSwitch(cmd, nonFlagArgs)
			// Then run namespace selector
			nsCmd.Run(nsCmd, []string{})
		} else {
			// No namespace argument found, use normal flow
			runSwitch(cmd, nonFlagArgs)
		}
	},
}

// contextCompletion provides autocompletion for context names
func contextCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	contexts, err := kubeconfig.GetContexts()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var suggestions []string
	for name := range contexts {
		suggestions = append(suggestions, name)
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(switchCmd)

	// Add flags
	switchCmd.Flags().BoolP("set-namespace", "n", false, "Also set the namespace after switching context")
}
