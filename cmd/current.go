package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/user-cube/kontext/pkg/kubeconfig"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current Kubernetes context",
	Long:  `Display the currently active Kubernetes context from your kubeconfig file.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentContext, err := kubeconfig.GetCurrentContext()
		if err != nil {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s Error retrieving current context: %v\n", red("✗"), err)
			os.Exit(1)
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		green := color.New(color.FgGreen, color.Bold).SprintFunc()

		fmt.Printf("%s %s %s\n", green("→"), bold("Current context:"), cyan(currentContext))
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
