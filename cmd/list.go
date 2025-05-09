package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Kubernetes contexts",
	Long:  `List all available Kubernetes contexts from your kubeconfig file.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		// Create colored output
		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		faint := color.New(color.Faint).SprintFunc()

		// Print header
		fmt.Println(bold("Available Kubernetes contexts:"))
		fmt.Println(faint("───────────────────────────────────"))

		// Sort context names for consistent display
		contextNames := make([]string, 0, len(contexts))
		for name := range contexts {
			contextNames = append(contextNames, name)
		}
		sort.Strings(contextNames)

		// Print contexts
		for _, name := range contextNames {
			if name == currentContext {
				fmt.Printf("%s %s %s\n", green("→"), cyan(name), green("(current)"))
			} else {
				fmt.Printf("  %s\n", name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
