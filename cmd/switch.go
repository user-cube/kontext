package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
)

// runSwitch contains the main logic for the switch command
// This is shared with the root command to enable the same behavior with `kontext`
func runSwitch(cmd *cobra.Command, args []string) {
	var contextName string

	// Get available contexts
	contexts, err := kubeconfig.GetContexts()
	if err != nil {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s Error retrieving contexts: %v\n", red("✗"), err)
		os.Exit(1)
	}

	currentContext, err := kubeconfig.GetCurrentContext()
	if err != nil {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s Error retrieving current context: %v\n", red("✗"), err)
		os.Exit(1)
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

		// Sort context names for consistent display
		sort.Strings(contextNames)

		// Check if we should set namespace after context switch
		setNS, _ := cmd.Flags().GetBool("set-namespace")

		// Prepare selector with enhanced templates and styling
		templates := &promptui.SelectTemplates{
			Label:    "{{ \"Select Kubernetes Context:\" | bold }}",
			Active:   "{{ \"→\" | cyan | bold }} {{ . | cyan | bold }}{{ if eq . \"" + currentContext + "\" }} {{ \"(current)\" | green | bold }}{{ end }}",
			Inactive: "  {{ . }}{{ if eq . \"" + currentContext + "\" }} {{ \"(current)\" | green }}{{ end }}",
			Selected: "{{ \"✓\" | green | bold }} {{ \"Selected context:\" | bold }} {{ . | cyan | bold }}",
			Details:  "{{ \"───────────────────────────────────────\" | faint }}\n{{ \"  Use arrow keys to navigate and Enter to select\" | faint }}",
		}

		prompt := promptui.Select{
			Label:     "Context",
			Items:     contextNames,
			Templates: templates,
			Size:      10,
			// Start with the current context selected
			CursorPos: getContextPosition(contextNames, currentContext),
		}

		_, selection, err := prompt.Run()
		if err != nil {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s Selection canceled: %v\n", red("✗"), err)
			return
		}

		contextName = selection

		// Don't switch if selected context is already current
		if contextName == currentContext {
			yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
			fmt.Printf("%s Context '%s' is already selected\n", yellow("!"), contextName)

			// Show current namespace info
			cyan := color.New(color.FgCyan).SprintFunc()
			green := color.New(color.FgGreen, color.Bold).SprintFunc()
			fmt.Printf("%s Current namespace: %s\n", green("→"), cyan(currentNamespace))
			return
		}

		// Switch to the selected context
		err = kubeconfig.SwitchContext(contextName)
		if err != nil {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s Error switching context: %v\n", red("✗"), err)
			os.Exit(1)
		}

		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Printf("%s Switched to context %s\n", green("✓"), cyan(contextName))

		// Get namespace for the new context
		targetNamespace, err := kubeconfig.GetNamespaceForContext(contextName)
		if err != nil {
			// Non-fatal - default to "default" namespace
			targetNamespace = "default"
		}

		fmt.Printf("%s Namespace: %s\n", green("→"), cyan(targetNamespace))

		// If requested, also let the user select a namespace
		if setNS {
			// We're already in the new context
			nsCmd.Run(nsCmd, []string{})
		}
	} else {
		contextName = args[0]

		// Don't switch if selected context is already current
		if contextName == currentContext {
			yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
			fmt.Printf("%s Context '%s' is already selected\n", yellow("!"), contextName)

			// Show current namespace info
			cyan := color.New(color.FgCyan).SprintFunc()
			green := color.New(color.FgGreen, color.Bold).SprintFunc()
			fmt.Printf("%s Current namespace: %s\n", green("→"), cyan(currentNamespace))
			return
		}

		// Switch to the selected context
		err = kubeconfig.SwitchContext(contextName)
		if err != nil {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s Error switching context: %v\n", red("✗"), err)
			os.Exit(1)
		}

		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Printf("%s Switched to context %s\n", green("✓"), cyan(contextName))

		// Get namespace for the new context
		targetNamespace, err := kubeconfig.GetNamespaceForContext(contextName)
		if err != nil {
			// Non-fatal - default to "default" namespace
			targetNamespace = "default"
		}

		fmt.Printf("%s Namespace: %s\n", green("→"), cyan(targetNamespace))

		// If requested, also let the user select a namespace
		setNS, _ := cmd.Flags().GetBool("set-namespace")
		if setNS {
			// We're already in the new context
			nsCmd.Run(nsCmd, []string{})
		}
	}
}

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [context]",
	Short: "Switch to a specific Kubernetes context",
	Long: `Switch to a specific Kubernetes context in your kubeconfig file.
If no context is provided, an interactive selection menu will be displayed.`,
	ValidArgsFunction: contextCompletion,
	Run:               runSwitch,
}

// Helper function to find the position of the current context in the list
func getContextPosition(contexts []string, currentContext string) int {
	for i, ctx := range contexts {
		if ctx == currentContext {
			return i
		}
	}
	return 0
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
