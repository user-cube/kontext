package kubeconfig

import (
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// createTestKubeConfig creates a temporary kubeconfig file for testing
func createTestKubeConfig(t *testing.T) (string, *api.Config) {
	t.Helper()

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "kontext-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create a test kubeconfig
	config := api.NewConfig()
	config.Clusters["cluster1"] = &api.Cluster{
		Server: "https://cluster1.example.com",
	}
	config.Clusters["cluster2"] = &api.Cluster{
		Server: "https://cluster2.example.com",
	}
	config.AuthInfos["user1"] = &api.AuthInfo{
		Token: "token1",
	}
	config.AuthInfos["user2"] = &api.AuthInfo{
		Token: "token2",
	}
	config.Contexts["context1"] = &api.Context{
		Cluster:   "cluster1",
		AuthInfo:  "user1",
		Namespace: "namespace1",
	}
	config.Contexts["context2"] = &api.Context{
		Cluster:   "cluster2",
		AuthInfo:  "user2",
		Namespace: "",
	}
	config.Contexts["context3"] = &api.Context{
		Cluster:   "cluster1", // Shared cluster
		AuthInfo:  "user1",    // Shared authInfo
		Namespace: "namespace3",
	}
	config.CurrentContext = "context1"

	configPath := filepath.Join(tmpDir, "config")
	if err := clientcmd.WriteToFile(*config, configPath); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	return configPath, config
}

func TestGetKubeConfigPath(t *testing.T) {
	tests := []struct {
		name           string
		kubeconfigEnv  string
		expectedPrefix string
	}{
		{
			name:           "Uses KUBECONFIG env var",
			kubeconfigEnv:  "/custom/path/config",
			expectedPrefix: "/custom/path/config",
		},
		{
			name:           "Falls back to default",
			kubeconfigEnv:  "",
			expectedPrefix: filepath.Join(os.Getenv("HOME"), ".kube", "config"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalEnv := os.Getenv("KUBECONFIG")
			defer func() {
				_ = os.Setenv("KUBECONFIG", originalEnv)
			}()

			// Set test environment
			if tt.kubeconfigEnv != "" {
				_ = os.Setenv("KUBECONFIG", tt.kubeconfigEnv)
			} else {
				_ = os.Unsetenv("KUBECONFIG")
			}

			path := GetKubeConfigPath()
			if path != tt.expectedPrefix {
				t.Errorf("GetKubeConfigPath() = %v, want %v", path, tt.expectedPrefix)
			}
		})
	}
}

func TestGetKubeConfig(t *testing.T) {
	configPath, expectedConfig := createTestKubeConfig(t)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(configPath))
	}()

	// Save original env
	originalEnv := os.Getenv("KUBECONFIG")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalEnv)
	}()

	_ = os.Setenv("KUBECONFIG", configPath)

	config, err := GetKubeConfig()
	if err != nil {
		t.Fatalf("GetKubeConfig() error = %v", err)
	}

	if config == nil {
		t.Fatal("GetKubeConfig() returned nil config")
	}

	// Verify contexts
	if len(config.Contexts) != len(expectedConfig.Contexts) {
		t.Errorf("GetKubeConfig() contexts count = %d, want %d", len(config.Contexts), len(expectedConfig.Contexts))
	}

	// Verify current context
	if config.CurrentContext != expectedConfig.CurrentContext {
		t.Errorf("GetKubeConfig() current context = %v, want %v", config.CurrentContext, expectedConfig.CurrentContext)
	}
}

func TestGetContexts(t *testing.T) {
	configPath, expectedConfig := createTestKubeConfig(t)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(configPath))
	}()

	originalEnv := os.Getenv("KUBECONFIG")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalEnv)
	}()

	_ = os.Setenv("KUBECONFIG", configPath)

	contexts, err := GetContexts()
	if err != nil {
		t.Fatalf("GetContexts() error = %v", err)
	}

	if len(contexts) != len(expectedConfig.Contexts) {
		t.Errorf("GetContexts() count = %d, want %d", len(contexts), len(expectedConfig.Contexts))
	}

	// Verify all expected contexts exist
	for name := range expectedConfig.Contexts {
		if _, exists := contexts[name]; !exists {
			t.Errorf("GetContexts() missing context: %s", name)
		}
	}
}

func TestGetCurrentContext(t *testing.T) {
	configPath, expectedConfig := createTestKubeConfig(t)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(configPath))
	}()

	originalEnv := os.Getenv("KUBECONFIG")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalEnv)
	}()

	_ = os.Setenv("KUBECONFIG", configPath)

	current, err := GetCurrentContext()
	if err != nil {
		t.Fatalf("GetCurrentContext() error = %v", err)
	}

	if current != expectedConfig.CurrentContext {
		t.Errorf("GetCurrentContext() = %v, want %v", current, expectedConfig.CurrentContext)
	}
}

func TestSwitchContext(t *testing.T) {
	configPath, _ := createTestKubeConfig(t)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(configPath))
	}()

	originalEnv := os.Getenv("KUBECONFIG")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalEnv)
	}()

	_ = os.Setenv("KUBECONFIG", configPath)

	tests := []struct {
		name        string
		contextName string
		wantErr     bool
		wantCurrent string
	}{
		{
			name:        "Switch to existing context",
			contextName: "context2",
			wantErr:     false,
			wantCurrent: "context2",
		},
		{
			name:        "Switch to non-existent context",
			contextName: "nonexistent",
			wantErr:     true,
		},
		{
			name:        "Switch to same context",
			contextName: "context1",
			wantErr:     false,
			wantCurrent: "context1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SwitchContext(tt.contextName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SwitchContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				current, err := GetCurrentContext()
				if err != nil {
					t.Fatalf("GetCurrentContext() error = %v", err)
				}
				if current != tt.wantCurrent {
					t.Errorf("GetCurrentContext() = %v, want %v", current, tt.wantCurrent)
				}
			}
		})
	}
}

func TestDeleteContext(t *testing.T) {
	tests := []struct {
		name               string
		contextToDelete    string
		initialCurrent     string
		wantErr            bool
		wantContextsCount  int
		wantCurrent        string
		wantClustersCount  int
		wantAuthInfosCount int
	}{
		{
			name:               "Delete non-current context",
			contextToDelete:    "context2",
			initialCurrent:     "context1",
			wantErr:            false,
			wantContextsCount:  2,
			wantCurrent:        "context1",
			wantClustersCount:  1, // cluster2 removed (only used by context2), cluster1 still used by context1 and context3
			wantAuthInfosCount: 1, // user2 removed (only used by context2), user1 still used by context1 and context3
		},
		{
			name:               "Delete current context",
			contextToDelete:    "context1",
			initialCurrent:     "context1",
			wantErr:            false,
			wantContextsCount:  2,
			wantCurrent:        "", // Current should be unset
			wantClustersCount:  2,  // cluster1 still used by context3, cluster2 still used by context2
			wantAuthInfosCount: 2,  // user1 still used by context3, user2 still used by context2
		},
		{
			name:               "Delete context with shared cluster and authInfo",
			contextToDelete:    "context3",
			initialCurrent:     "context1",
			wantErr:            false,
			wantContextsCount:  2,
			wantCurrent:        "context1",
			wantClustersCount:  2, // cluster1 still used by context1
			wantAuthInfosCount: 2, // user1 still used by context1, user2 still used by context2
		},
		{
			name:            "Delete non-existent context",
			contextToDelete: "nonexistent",
			initialCurrent:  "context1",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath, _ := createTestKubeConfig(t)
			defer func() {
				_ = os.RemoveAll(filepath.Dir(configPath))
			}()

			originalEnv := os.Getenv("KUBECONFIG")
			defer func() {
				_ = os.Setenv("KUBECONFIG", originalEnv)
			}()

			_ = os.Setenv("KUBECONFIG", configPath)

			// Set initial current context
			if err := SwitchContext(tt.initialCurrent); err != nil {
				t.Fatalf("Failed to set initial context: %v", err)
			}

			err := DeleteContext(tt.contextToDelete)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify contexts count
				contexts, err := GetContexts()
				if err != nil {
					t.Fatalf("GetContexts() error = %v", err)
				}
				if len(contexts) != tt.wantContextsCount {
					t.Errorf("GetContexts() count = %d, want %d", len(contexts), tt.wantContextsCount)
				}

				// Verify context was deleted
				if _, exists := contexts[tt.contextToDelete]; exists {
					t.Errorf("DeleteContext() context %s still exists", tt.contextToDelete)
				}

				// Verify current context
				current, err := GetCurrentContext()
				if err != nil && tt.wantCurrent != "" {
					t.Fatalf("GetCurrentContext() error = %v", err)
				}
				if current != tt.wantCurrent {
					t.Errorf("GetCurrentContext() = %v, want %v", current, tt.wantCurrent)
				}

				// Verify clusters and authInfos cleanup
				config, err := GetKubeConfig()
				if err != nil {
					t.Fatalf("GetKubeConfig() error = %v", err)
				}
				if len(config.Clusters) != tt.wantClustersCount {
					t.Errorf("Clusters count = %d, want %d", len(config.Clusters), tt.wantClustersCount)
				}
				if len(config.AuthInfos) != tt.wantAuthInfosCount {
					t.Errorf("AuthInfos count = %d, want %d", len(config.AuthInfos), tt.wantAuthInfosCount)
				}
			}
		})
	}
}

func TestGetCurrentNamespace(t *testing.T) {
	configPath, _ := createTestKubeConfig(t)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(configPath))
	}()

	originalEnv := os.Getenv("KUBECONFIG")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalEnv)
	}()

	_ = os.Setenv("KUBECONFIG", configPath)

	tests := []struct {
		name        string
		contextName string
		want        string
		wantErr     bool
	}{
		{
			name:        "Get namespace from context with namespace",
			contextName: "context1",
			want:        "namespace1",
			wantErr:     false,
		},
		{
			name:        "Get namespace from context without namespace (default)",
			contextName: "context2",
			want:        "default",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SwitchContext(tt.contextName); err != nil {
				t.Fatalf("SwitchContext() error = %v", err)
			}

			namespace, err := GetCurrentNamespace()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if namespace != tt.want {
				t.Errorf("GetCurrentNamespace() = %v, want %v", namespace, tt.want)
			}
		})
	}
}

func TestGetNamespaceForContext(t *testing.T) {
	configPath, _ := createTestKubeConfig(t)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(configPath))
	}()

	originalEnv := os.Getenv("KUBECONFIG")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalEnv)
	}()

	_ = os.Setenv("KUBECONFIG", configPath)

	tests := []struct {
		name        string
		contextName string
		want        string
		wantErr     bool
	}{
		{
			name:        "Get namespace for context with namespace",
			contextName: "context1",
			want:        "namespace1",
			wantErr:     false,
		},
		{
			name:        "Get namespace for context without namespace (default)",
			contextName: "context2",
			want:        "default",
			wantErr:     false,
		},
		{
			name:        "Get namespace for non-existent context",
			contextName: "nonexistent",
			want:        "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespace, err := GetNamespaceForContext(tt.contextName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNamespaceForContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if namespace != tt.want {
				t.Errorf("GetNamespaceForContext() = %v, want %v", namespace, tt.want)
			}
		})
	}
}

func TestSetNamespace(t *testing.T) {
	configPath, _ := createTestKubeConfig(t)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(configPath))
	}()

	originalEnv := os.Getenv("KUBECONFIG")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalEnv)
	}()

	_ = os.Setenv("KUBECONFIG", configPath)

	// Set current context
	if err := SwitchContext("context1"); err != nil {
		t.Fatalf("SwitchContext() error = %v", err)
	}

	tests := []struct {
		name      string
		namespace string
		wantErr   bool
		want      string
	}{
		{
			name:      "Set namespace for current context",
			namespace: "new-namespace",
			wantErr:   false,
			want:      "new-namespace",
		},
		{
			name:      "Set empty namespace (should work)",
			namespace: "",
			wantErr:   false,
			want:      "default", // GetCurrentNamespace returns "default" for empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetNamespace(tt.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := GetCurrentNamespace()
				if err != nil {
					t.Fatalf("GetCurrentNamespace() error = %v", err)
				}
				if got != tt.want {
					t.Errorf("GetCurrentNamespace() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSetNamespaceForContext(t *testing.T) {
	configPath, _ := createTestKubeConfig(t)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(configPath))
	}()

	originalEnv := os.Getenv("KUBECONFIG")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalEnv)
	}()

	_ = os.Setenv("KUBECONFIG", configPath)

	tests := []struct {
		name        string
		contextName string
		namespace   string
		wantErr     bool
		want        string
	}{
		{
			name:        "Set namespace for specific context",
			contextName: "context2",
			namespace:   "new-namespace",
			wantErr:     false,
			want:        "new-namespace",
		},
		{
			name:        "Set namespace for non-existent context",
			contextName: "nonexistent",
			namespace:   "new-namespace",
			wantErr:     true,
		},
		{
			name:        "Set namespace using empty context name (uses current)",
			contextName: "",
			namespace:   "current-namespace",
			wantErr:     false,
			want:        "current-namespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set current context for empty contextName test
			if tt.contextName == "" {
				if err := SwitchContext("context1"); err != nil {
					t.Fatalf("SwitchContext() error = %v", err)
				}
			}

			err := SetNamespaceForContext(tt.contextName, tt.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetNamespaceForContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				contextToCheck := tt.contextName
				if contextToCheck == "" {
					contextToCheck = "context1" // We set this as current
				}
				got, err := GetNamespaceForContext(contextToCheck)
				if err != nil {
					t.Fatalf("GetNamespaceForContext() error = %v", err)
				}
				// For empty namespace, GetNamespaceForContext returns "default"
				expected := tt.want
				if tt.namespace == "" {
					expected = "default"
				}
				if got != expected {
					t.Errorf("GetNamespaceForContext() = %v, want %v", got, expected)
				}
			}
		})
	}
}
