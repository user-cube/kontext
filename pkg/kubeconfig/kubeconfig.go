package kubeconfig

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// GetKubeConfigPath returns the path to the kubeconfig file
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
		return []string{"default", "kube-system", "kube-public", "kube-node-lease"}, nil
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return []string{"default", "kube-system", "kube-public", "kube-node-lease"}, nil
	}

	// Try to list namespaces from the cluster
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		// Fall back to default namespaces if there's an error
		return []string{"default", "kube-system", "kube-public", "kube-node-lease"}, nil
	}

	// Extract namespace names from the response
	namespaces := make([]string, 0, len(namespaceList.Items))
	for _, ns := range namespaceList.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}
