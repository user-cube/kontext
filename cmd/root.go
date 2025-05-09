/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
	"github.com/user-cube/kontext/pkg/ui"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kontext [context] [namespace]",
	Short: "A tool for managing Kubernetes contexts",
	Long: `Kontext is a CLI tool that helps you manage your Kubernetes contexts.
It allows you to list, view, and switch between different Kubernetes contexts
in your kubeconfig file with ease and tab completion support.

Examples:
  kontext list                      # List all available contexts
  kontext current                   # Show the current context
  kontext switch my-context         # Switch to a specific context
  kontext my-context                # Switch to a specific context
  kontext -n                        # Switch context and then select namespace
  kontext my-context -n             # Switch to context and then select namespace
  kontext my-context -n my-namespace # Switch to context and set namespace directly`,
	// When no subcommands are provided, run the switch command functionality
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Enable shell completion
	rootCmd.CompletionOptions.DisableDefaultCmd = false

	// Add flags - same as switch command
	rootCmd.Flags().BoolP("set-namespace", "n", false, "Also set the namespace after switching context")
}
