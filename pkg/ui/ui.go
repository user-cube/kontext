// Package ui provides UI utilities for the kontext CLI
//
// This package centralizes all user interface presentation logic, including:
// - Colorized output (success, errors, warnings, info)
// - Interactive selectors for contexts and namespaces
// - Consistent formatting and styling
package ui

import (
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// Colors creates and returns commonly used colored print functions
type Colors struct {
	Red    func(a ...interface{}) string
	Green  func(a ...interface{}) string
	Yellow func(a ...interface{}) string
	Cyan   func(a ...interface{}) string
	Bold   func(a ...interface{}) string
	Faint  func(a ...interface{}) string
}

// NewColors returns initialized color functions for consistent UI styling
func NewColors() *Colors {
	return &Colors{
		Red:    color.New(color.FgRed, color.Bold).SprintFunc(),
		Green:  color.New(color.FgGreen, color.Bold).SprintFunc(),
		Yellow: color.New(color.FgYellow, color.Bold).SprintFunc(),
		Cyan:   color.New(color.FgCyan, color.Bold).SprintFunc(),
		Bold:   color.New(color.Bold).SprintFunc(),
		Faint:  color.New(color.Faint).SprintFunc(),
	}
}

// PrintError prints a formatted error message and exits if exitOnError is true
// If err is nil, only the message is displayed
func PrintError(msg string, err error, exitOnError bool) {
	colors := NewColors()
	if err != nil {
		fmt.Printf("%s %s: %v\n", colors.Red("✗"), msg, err)
	} else {
		fmt.Printf("%s %s\n", colors.Red("✗"), msg)
	}
	if exitOnError {
		os.Exit(1)
	}
}

// PrintSuccess prints a formatted success message
func PrintSuccess(msg string, details ...string) {
	colors := NewColors()
	fmt.Printf("%s %s", colors.Green("✓"), msg)

	for _, detail := range details {
		fmt.Printf(" %s", colors.Cyan(detail))
	}
	fmt.Println()
}

// PrintWarning prints a formatted warning message
func PrintWarning(msg string, details ...string) {
	colors := NewColors()
	fmt.Printf("%s %s", colors.Yellow("!"), msg)

	for _, detail := range details {
		fmt.Printf(" %s", colors.Cyan(detail))
	}
	fmt.Println()
}

// PrintInfo prints a formatted information label and value
func PrintInfo(label string, value string) {
	colors := NewColors()
	fmt.Printf("%s: %s\n", colors.Bold(label), value)
}

// PrintNote prints a formatted note message with an info icon
func PrintNote(msg string, details ...string) {
	colors := NewColors()
	// Using blue color with info icon for notes
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	fmt.Printf("%s %s", blue("ℹ"), blue("Note:"))

	fmt.Printf(" %s", msg)

	for _, detail := range details {
		fmt.Printf(" %s", colors.Cyan(detail))
	}
	fmt.Println()
}

// PrintCurrentContext displays the current context information
func PrintCurrentContext(contextName string) {
	colors := NewColors()
	fmt.Printf("%s %s %s\n", colors.Green("→"), colors.Bold("Current context:"), colors.Cyan(contextName))
}

// PrintCurrentNamespace displays the current namespace information for a context
func PrintCurrentNamespace(contextName, namespaceName string) {
	colors := NewColors()
	fmt.Printf("%s %s %s\n", colors.Green("→"), colors.Bold(fmt.Sprintf("Context: %s Namespace:", contextName)), colors.Cyan(namespaceName))
}

// CreateContextSelector creates an interactive prompt UI for selecting Kubernetes contexts
func CreateContextSelector(contexts []string, currentContext string) *promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "{{ \"Select Kubernetes Context:\" | bold }}",
		Active:   "{{ \"→\" | cyan | bold }} {{ . | cyan | bold }}{{ if eq . \"" + currentContext + "\" }} {{ \"(current)\" | green | bold }}{{ end }}",
		Inactive: "  {{ . }}{{ if eq . \"" + currentContext + "\" }} {{ \"(current)\" | green }}{{ end }}",
		Selected: "{{ \"✓\" | green | bold }} {{ \"Selected context:\" | bold }} {{ . | cyan | bold }}",
		Details:  "{{ \"───────────────────────────────────────\" | faint }}\n{{ \"  Use arrow keys to navigate and Enter to select\" | faint }}",
	}

	cursorPos := 0
	for i, ctx := range contexts {
		if ctx == currentContext {
			cursorPos = i
			break
		}
	}

	// Log statement to help debugging cursor position issues
	if cursorPos == 0 && currentContext != "" && len(contexts) > 0 && contexts[0] != currentContext {
		PrintNote(fmt.Sprintf("Current context '%s' not found in sorted context list, defaulting to first item", currentContext))
	}

	return &promptui.Select{
		Label:     "Context",
		Items:     contexts,
		Templates: templates,
		Size:      10,
		CursorPos: cursorPos,
	}
}

// CreateNamespaceSelector creates an interactive prompt UI for selecting Kubernetes namespaces
func CreateNamespaceSelector(namespaces []string, currentNamespace string, currentContext string) *promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "{{ \"Select Namespace:\" | bold }}",
		Active:   "{{ \"→\" | cyan | bold }} {{ . | cyan | bold }}{{ if eq . \"" + currentNamespace + "\" }} {{ \"(current)\" | green | bold }}{{ end }}",
		Inactive: "  {{ . }}{{ if eq . \"" + currentNamespace + "\" }} {{ \"(current)\" | green }}{{ end }}",
		Selected: "{{ \"✓\" | green | bold }} {{ \"Context:\" | bold }} {{ \"" + currentContext + "\" | cyan | bold }} {{ \"Namespace:\" | bold }} {{ . | cyan | bold }}",
		Details:  "{{ \"───────────────────────────────────────\" | faint }}\n{{ \"  Use arrow keys to navigate and Enter to select\" | faint }}",
	}

	cursorPos := 0
	for i, ns := range namespaces {
		if ns == currentNamespace {
			cursorPos = i
			break
		}
	}

	// Log statement to help debugging cursor position issues
	if cursorPos == 0 && currentNamespace != "" && len(namespaces) > 0 && namespaces[0] != currentNamespace {
		PrintNote(fmt.Sprintf("Current namespace '%s' not found in sorted namespace list, defaulting to first item", currentNamespace))
	}

	return &promptui.Select{
		Label:     "Namespace",
		Items:     namespaces,
		Templates: templates,
		Size:      10,
		CursorPos: cursorPos,
	}
}

// SortContexts sorts the context names alphabetically, optionally placing the current context first
// This function can be used to ensure the current context is always at the top of the list
// If prioritizeCurrent is true, the current context will be placed first in the sorted list
func SortContexts(contextNames []string, currentContext string, prioritizeCurrent bool) []string {
	// Handle empty context list
	if len(contextNames) == 0 {
		return contextNames
	}

	// Make a copy of the slice to avoid modifying the original
	sorted := make([]string, len(contextNames))
	copy(sorted, contextNames)

	// Sort alphabetically first
	sort.Strings(sorted)

	// If requested, move current context to the front
	if prioritizeCurrent && currentContext != "" {
		found := false
		// Find current context and move to front
		for i, ctx := range sorted {
			if ctx == currentContext {
				// Remove the current context from its position
				sorted = append(sorted[:i], sorted[i+1:]...)
				// Add it to the front
				sorted = append([]string{currentContext}, sorted...)
				found = true
				break
			}
		} // If current context wasn't found in the list, log a debug message
		if !found {
			PrintNote(fmt.Sprintf("Current context '%s' not found in context list", currentContext))
		}
	}

	return sorted
}

// SortNamespaces sorts the namespace names alphabetically, optionally placing the current namespace first
// This function can be used to ensure the current namespace is always at the top of the list
// If prioritizeCurrent is true, the current namespace will be placed first in the sorted list
func SortNamespaces(namespaces []string, currentNamespace string, prioritizeCurrent bool) []string {
	// Handle empty namespace list
	if len(namespaces) == 0 {
		return namespaces
	}

	// Make a copy of the slice to avoid modifying the original
	sorted := make([]string, len(namespaces))
	copy(sorted, namespaces)

	// Sort alphabetically first
	sort.Strings(sorted)

	// If requested, move current namespace to the front
	if prioritizeCurrent && currentNamespace != "" {
		found := false
		// Find current namespace and move to front
		for i, ns := range sorted {
			if ns == currentNamespace {
				// Remove the current namespace from its position
				sorted = append(sorted[:i], sorted[i+1:]...)
				// Add it to the front
				sorted = append([]string{currentNamespace}, sorted...)
				found = true
				break
			}
		} // If current namespace wasn't found in the list, log a debug message
		if !found {
			PrintNote(fmt.Sprintf("Current namespace '%s' not found in namespace list", currentNamespace))
		}
	}

	return sorted
}

// PrintContextList displays a formatted list of available contexts
// Output:
//
//	Available Kubernetes contexts:
//	───────────────────────────────────
//	dev
//	staging
//	→ production (current)
func PrintContextList(contextNames []string, currentContext string) {
	colors := NewColors()

	// Print header
	fmt.Println(colors.Bold("Available Kubernetes contexts:"))
	fmt.Println(colors.Faint("───────────────────────────────────"))

	// Print contexts
	for _, name := range contextNames {
		if name == currentContext {
			fmt.Printf("%s %s %s\n", colors.Green("→"), colors.Cyan(name), colors.Green("(current)"))
		} else {
			fmt.Printf("  %s\n", name)
		}
	}
}

// ConfirmAction shows a yes/no confirmation prompt for potentially destructive actions.
// It returns true if the user confirms, false if the user cancels.
func ConfirmAction(message string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		// When IsConfirm is true, promptui returns ErrAbort on "no"/cancel.
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
