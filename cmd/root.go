/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kontext",
	Short: "A tool for managing Kubernetes contexts",
	Long: `Kontext is a CLI tool that helps you manage your Kubernetes contexts.
It allows you to list, view, and switch between different Kubernetes contexts
in your kubeconfig file with ease and tab completion support.

Examples:
  kontext list                # List all available contexts
  kontext current             # Show the current context
  kontext switch my-context   # Switch to a specific context`,
	// When no subcommands are provided, run the switch command functionality
	Run: runSwitch,
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
}
