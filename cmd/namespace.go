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

// nsCmd represents the namespace command
var nsCmd = &cobra.Command{
	Use:     "namespace [namespace]",
	Aliases: []string{"ns"},
	Short:   "View or change the current namespace",
	Long: `View or change the namespace for the current Kubernetes context.
If no namespace is provided, an interactive selection menu will be displayed.
Examples:
  kontext namespace                  # Show interactive selector for namespaces
  kontext namespace my-namespace     # Switch to the specified namespace`,
	Run: runNamespace,
}

// Function to run the namespace command
func runNamespace(cmd *cobra.Command, args []string) {
	// Get current context and namespace
	currentContext, err := kubeconfig.GetCurrentContext()
	if err != nil {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s Error retrieving current context: %v\n", red("✗"), err)
		os.Exit(1)
	}

	currentNamespace, err := kubeconfig.GetCurrentNamespace()
	if err != nil {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s Error retrieving current namespace: %v\n", red("✗"), err)
		os.Exit(1)
	}

	// If no arguments are provided, show the current namespace or interactive selector
	if len(args) == 0 {
		contextStr := fmt.Sprintf("Context: %s", currentContext)
		namespaceStr := fmt.Sprintf("Namespace: %s", currentNamespace)

		// Check if a specific flag was provided to only show the current value
		if showOnly, _ := cmd.Flags().GetBool("show"); showOnly {
			bold := color.New(color.Bold).SprintFunc()
			cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
			green := color.New(color.FgGreen, color.Bold).SprintFunc()

			fmt.Printf("%s %s %s\n", green("→"), bold(contextStr), cyan(namespaceStr))
			return
		}

		// Get available namespaces
		namespaces, err := kubeconfig.GetNamespaces()
		if err != nil {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s Error retrieving namespaces: %v\n", red("✗"), err)
			os.Exit(1)
		}

		// Sort namespaces for consistent display
		sort.Strings(namespaces)

		// Create an interactive selector
		templates := &promptui.SelectTemplates{
			Label:    "{{ \"Select Namespace:\" | bold }}",
			Active:   "{{ \"→\" | cyan | bold }} {{ . | cyan | bold }}{{ if eq . \"" + currentNamespace + "\" }} {{ \"(current)\" | green | bold }}{{ end }}",
			Inactive: "  {{ . }}{{ if eq . \"" + currentNamespace + "\" }} {{ \"(current)\" | green }}{{ end }}",
			Selected: "{{ \"✓\" | green | bold }} {{ \"Context:\" | bold }} {{ \"" + currentContext + "\" | cyan | bold }} {{ \"Namespace:\" | bold }} {{ . | cyan | bold }}",
			Details:  "{{ \"───────────────────────────────────────\" | faint }}\n{{ \"  Use arrow keys to navigate and Enter to select\" | faint }}",
		}

		prompt := promptui.Select{
			Label:     "Namespace",
			Items:     namespaces,
			Templates: templates,
			Size:      10,
			// Start with the current namespace selected
			CursorPos: getNamespacePosition(namespaces, currentNamespace),
		}

		_, selection, err := prompt.Run()
		if err != nil {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s Selection canceled: %v\n", red("✗"), err)
			return
		}

		// If the selected namespace is the same as the current one, don't do anything
		if selection == currentNamespace {
			yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
			fmt.Printf("%s Namespace '%s' is already selected\n", yellow("!"), selection)
			return
		}

		// Change the namespace
		err = kubeconfig.SetNamespace(selection)
		if err != nil {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s Error setting namespace: %v\n", red("✗"), err)
			os.Exit(1)
		}

		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Printf("%s Switched to namespace %s in context %s\n", green("✓"), cyan(selection), cyan(currentContext))
		return
	}

	// Change to the specified namespace
	namespace := args[0]

	// If the specified namespace is the same as the current one, don't do anything
	if namespace == currentNamespace {
		yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
		fmt.Printf("%s Namespace '%s' is already selected\n", yellow("!"), namespace)
		return
	}

	// Change the namespace
	err = kubeconfig.SetNamespace(namespace)
	if err != nil {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s Error setting namespace: %v\n", red("✗"), err)
		os.Exit(1)
	}

	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s Switched to namespace %s in context %s\n", green("✓"), cyan(namespace), cyan(currentContext))
}

// Helper function to find the position of the current namespace in the list
func getNamespacePosition(namespaces []string, currentNamespace string) int {
	for i, ns := range namespaces {
		if ns == currentNamespace {
			return i
		}
	}
	return 0
}

func init() {
	rootCmd.AddCommand(nsCmd)

	// Add flags
	nsCmd.Flags().BoolP("show", "s", false, "Only show the current namespace without the selector")
}
