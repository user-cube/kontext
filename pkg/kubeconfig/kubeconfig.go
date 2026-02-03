// Package kubeconfig provides functions for interacting with Kubernetes configuration
//
// This package offers a simplified interface for common operations related to
// Kubernetes contexts and namespaces, including:
// - Retrieving available contexts
// - Getting and setting the current context
// - Managing namespaces
// - Connecting to clusters
package kubeconfig

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user-cube/kontext/pkg/static"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// GetKubeConfigPath returns the path to the kubeconfig file
//
// This function checks the KUBECONFIG environment variable first,
// and falls back to the default ~/.kube/config location if not set.
func GetKubeConfigPath() string {
	if os.Getenv("KUBECONFIG") != "" {
		return os.Getenv("KUBECONFIG")
	}
	return filepath.Join(os.Getenv("HOME"), ".kube", "config")
}

// GetKubeConfig loads the kubeconfig file
func GetKubeConfig() (*api.Config, error) {
	configPath := GetKubeConfigPath()
	config, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %w", err)
	}
	return config, nil
}

// GetContexts returns all available contexts in the kubeconfig
func GetContexts() (map[string]*api.Context, error) {
	config, err := GetKubeConfig()
	if err != nil {
		return nil, err
	}
	return config.Contexts, nil
}

// GetCurrentContext returns the name of the current context
func GetCurrentContext() (string, error) {
	config, err := GetKubeConfig()
	if err != nil {
		return "", err
	}
	return config.CurrentContext, nil
}

// SwitchContext changes the current context to the specified one
func SwitchContext(contextName string) error {
	configPath := GetKubeConfigPath()
	config, err := GetKubeConfig()
	if err != nil {
		return err
	}

	// Check if the context exists
	if _, exists := config.Contexts[contextName]; !exists {
		return fmt.Errorf("context '%s' does not exist", contextName)
	}

	// Set the current context
	config.CurrentContext = contextName

	// Save the updated config
	err = clientcmd.WriteToFile(*config, configPath)
	if err != nil {
		return fmt.Errorf("error saving kubeconfig: %w", err)
	}

	return nil
}

// DeleteContext removes the specified context from the kubeconfig.
// If the deleted context is the current context, the current context will be unset.
// Any clusters or authInfos that are no longer referenced by any remaining context
// will also be removed to keep the config clean.
func DeleteContext(contextName string) error {
	configPath := GetKubeConfigPath()
	config, err := GetKubeConfig()
	if err != nil {
		return err
	}

	// Check if the context exists
	ctx, exists := config.Contexts[contextName]
	if !exists {
		return fmt.Errorf("context '%s' does not exist", contextName)
	}

	// Track associated cluster and auth info so we can clean them up if unused
	clusterName := ctx.Cluster
	authInfoName := ctx.AuthInfo

	// Delete the context
	delete(config.Contexts, contextName)

	// Unset current context if it was the one being deleted
	if config.CurrentContext == contextName {
		config.CurrentContext = ""
	}

	// Helper to check if a cluster/authInfo is still referenced by any context
	isClusterReferenced := func(name string) bool {
		if name == "" {
			return false
		}
		for _, c := range config.Contexts {
			if c != nil && c.Cluster == name {
				return true
			}
		}
		return false
	}

	isAuthInfoReferenced := func(name string) bool {
		if name == "" {
			return false
		}
		for _, c := range config.Contexts {
			if c != nil && c.AuthInfo == name {
				return true
			}
		}
		return false
	}

	// Clean up cluster if no longer referenced
	if clusterName != "" && !isClusterReferenced(clusterName) {
		delete(config.Clusters, clusterName)
	}

	// Clean up authInfo if no longer referenced
	if authInfoName != "" && !isAuthInfoReferenced(authInfoName) {
		delete(config.AuthInfos, authInfoName)
	}

	// Save the updated config
	if err := clientcmd.WriteToFile(*config, configPath); err != nil {
		return fmt.Errorf("error saving kubeconfig: %w", err)
	}

	return nil
}

// GetCurrentNamespace returns the namespace set for the current context
func GetCurrentNamespace() (string, error) {
	config, err := GetKubeConfig()
	if err != nil {
		return "", err
	}

	currentContext := config.CurrentContext
	if currentContext == "" {
		return "", fmt.Errorf("no current context set")
	}

	context, exists := config.Contexts[currentContext]
	if !exists {
		return "", fmt.Errorf("current context '%s' does not exist in config", currentContext)
	}

	// If namespace is empty, return "default"
	if context.Namespace == "" {
		return "default", nil
	}

	return context.Namespace, nil
}

// GetNamespaceForContext returns the namespace for the specified context
func GetNamespaceForContext(contextName string) (string, error) {
	config, err := GetKubeConfig()
	if err != nil {
		return "", err
	}

	context, exists := config.Contexts[contextName]
	if !exists {
		return "", fmt.Errorf("context '%s' does not exist", contextName)
	}

	// If namespace is empty, return "default"
	if context.Namespace == "" {
		return "default", nil
	}

	return context.Namespace, nil
}

// SetNamespace sets the namespace for the current context
func SetNamespace(namespace string) error {
	return SetNamespaceForContext("", namespace)
}

// SetNamespaceForContext sets the namespace for the specified context
// If contextName is empty, it uses the current context
func SetNamespaceForContext(contextName string, namespace string) error {
	configPath := GetKubeConfigPath()
	config, err := GetKubeConfig()
	if err != nil {
		return err
	}

	// Use current context if none specified
	if contextName == "" {
		contextName = config.CurrentContext
		if contextName == "" {
			return fmt.Errorf("no current context set")
		}
	}

	// Check if the context exists
	context, exists := config.Contexts[contextName]
	if !exists {
		return fmt.Errorf("context '%s' does not exist", contextName)
	}

	// Set the namespace
	context.Namespace = namespace

	// Save the updated config
	err = clientcmd.WriteToFile(*config, configPath)
	if err != nil {
		return fmt.Errorf("error saving kubeconfig: %w", err)
	}

	return nil
}

// GetNamespaces returns all available namespaces for the current context
//
// This function attempts to connect to the cluster and list namespaces.
// If the connection fails, it returns a set of default namespaces.
func GetNamespaces() ([]string, error) {
	return GetNamespacesForContext("")
}

// GetNamespacesForContext returns all available namespaces for the specified context
// If contextName is empty, it uses the current context
func GetNamespacesForContext(contextName string) ([]string, error) {
	config, err := GetKubeConfig()
	if err != nil {
		return nil, err
	}

	// Use current context if none specified
	if contextName == "" {
		contextName = config.CurrentContext
		if contextName == "" {
			return nil, fmt.Errorf("no current context set")
		}
	}

	// Check if the context exists
	_, exists := config.Contexts[contextName]
	if !exists {
		return nil, fmt.Errorf("context '%s' does not exist", contextName)
	}

	// Create client configuration for the specified context
	clientConfig := clientcmd.NewNonInteractiveClientConfig(
		*config,
		contextName,
		&clientcmd.ConfigOverrides{},
		clientcmd.NewDefaultClientConfigLoadingRules(),
	)

	// Get REST config for the context
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		// If we can't connect to the cluster, return some default namespaces
		// This handles the case where the user might be offline or the cluster is unavailable
		return static.FallBackNamespace, nil
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return static.FallBackNamespace, nil
	}

	// Try to list namespaces from the cluster
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		// Fall back to default namespaces if there's an error
		return static.FallBackNamespace, nil
	}

	// Extract namespace names from the response
	namespaces := make([]string, 0, len(namespaceList.Items))
	for _, ns := range namespaceList.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}
