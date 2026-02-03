package ui

import (
	"testing"
)

func TestSortContexts(t *testing.T) {
	tests := []struct {
		name              string
		contexts          []string
		currentContext    string
		prioritizeCurrent bool
		want              []string
	}{
		{
			name:              "Sort alphabetically without prioritizing",
			contexts:          []string{"zebra", "alpha", "beta"},
			currentContext:    "beta",
			prioritizeCurrent: false,
			want:              []string{"alpha", "beta", "zebra"},
		},
		{
			name:              "Sort with current context prioritized",
			contexts:          []string{"zebra", "alpha", "beta"},
			currentContext:    "beta",
			prioritizeCurrent: true,
			want:              []string{"beta", "alpha", "zebra"},
		},
		{
			name:              "Empty list",
			contexts:          []string{},
			currentContext:    "",
			prioritizeCurrent: true,
			want:              []string{},
		},
		{
			name:              "Current context not in list",
			contexts:          []string{"alpha", "beta", "gamma"},
			currentContext:    "delta",
			prioritizeCurrent: true,
			want:              []string{"alpha", "beta", "gamma"}, // Should still sort alphabetically
		},
		{
			name:              "Single context",
			contexts:          []string{"single"},
			currentContext:    "single",
			prioritizeCurrent: true,
			want:              []string{"single"},
		},
		{
			name:              "Multiple contexts, current at end",
			contexts:          []string{"alpha", "beta", "gamma", "delta"},
			currentContext:    "delta",
			prioritizeCurrent: true,
			want:              []string{"delta", "alpha", "beta", "gamma"},
		},
		{
			name:              "Multiple contexts, current at beginning",
			contexts:          []string{"alpha", "beta", "gamma", "delta"},
			currentContext:    "alpha",
			prioritizeCurrent: true,
			want:              []string{"alpha", "beta", "delta", "gamma"},
		},
		{
			name:              "Empty current context",
			contexts:          []string{"alpha", "beta", "gamma"},
			currentContext:    "",
			prioritizeCurrent: true,
			want:              []string{"alpha", "beta", "gamma"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortContexts(tt.contexts, tt.currentContext, tt.prioritizeCurrent)
			if len(got) != len(tt.want) {
				t.Errorf("SortContexts() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("SortContexts()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSortNamespaces(t *testing.T) {
	tests := []struct {
		name              string
		namespaces        []string
		currentNamespace  string
		prioritizeCurrent bool
		want              []string
	}{
		{
			name:              "Sort alphabetically without prioritizing",
			namespaces:        []string{"zebra", "alpha", "beta"},
			currentNamespace:  "beta",
			prioritizeCurrent: false,
			want:              []string{"alpha", "beta", "zebra"},
		},
		{
			name:              "Sort with current namespace prioritized",
			namespaces:        []string{"zebra", "alpha", "beta"},
			currentNamespace:  "beta",
			prioritizeCurrent: true,
			want:              []string{"beta", "alpha", "zebra"},
		},
		{
			name:              "Empty list",
			namespaces:        []string{},
			currentNamespace:  "",
			prioritizeCurrent: true,
			want:              []string{},
		},
		{
			name:              "Current namespace not in list",
			namespaces:        []string{"alpha", "beta", "gamma"},
			currentNamespace:  "delta",
			prioritizeCurrent: true,
			want:              []string{"alpha", "beta", "gamma"}, // Should still sort alphabetically
		},
		{
			name:              "Single namespace",
			namespaces:        []string{"single"},
			currentNamespace:  "single",
			prioritizeCurrent: true,
			want:              []string{"single"},
		},
		{
			name:              "Multiple namespaces, current at end",
			namespaces:        []string{"alpha", "beta", "gamma", "delta"},
			currentNamespace:  "delta",
			prioritizeCurrent: true,
			want:              []string{"delta", "alpha", "beta", "gamma"},
		},
		{
			name:              "Multiple namespaces, current at beginning",
			namespaces:        []string{"alpha", "beta", "gamma", "delta"},
			currentNamespace:  "alpha",
			prioritizeCurrent: true,
			want:              []string{"alpha", "beta", "delta", "gamma"},
		},
		{
			name:              "Empty current namespace",
			namespaces:        []string{"alpha", "beta", "gamma"},
			currentNamespace:  "",
			prioritizeCurrent: true,
			want:              []string{"alpha", "beta", "gamma"},
		},
		{
			name:              "Kubernetes default namespaces",
			namespaces:        []string{"kube-system", "default", "kube-public", "kube-node-lease"},
			currentNamespace:  "default",
			prioritizeCurrent: true,
			want:              []string{"default", "kube-node-lease", "kube-public", "kube-system"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortNamespaces(tt.namespaces, tt.currentNamespace, tt.prioritizeCurrent)
			if len(got) != len(tt.want) {
				t.Errorf("SortNamespaces() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("SortNamespaces()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSortContextsDoesNotModifyOriginal(t *testing.T) {
	original := []string{"zebra", "alpha", "beta"}
	originalCopy := make([]string, len(original))
	copy(originalCopy, original)

	SortContexts(original, "beta", true)

	// Verify original slice is unchanged
	for i := range original {
		if original[i] != originalCopy[i] {
			t.Errorf("SortContexts() modified original slice at index %d: got %v, want %v", i, original[i], originalCopy[i])
		}
	}
}

func TestSortNamespacesDoesNotModifyOriginal(t *testing.T) {
	original := []string{"zebra", "alpha", "beta"}
	originalCopy := make([]string, len(original))
	copy(originalCopy, original)

	SortNamespaces(original, "beta", true)

	// Verify original slice is unchanged
	for i := range original {
		if original[i] != originalCopy[i] {
			t.Errorf("SortNamespaces() modified original slice at index %d: got %v, want %v", i, original[i], originalCopy[i])
		}
	}
}
