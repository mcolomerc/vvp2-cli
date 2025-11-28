package api

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNamespaceParsing(t *testing.T) {
	yamlContent := `
metadata:
  name: test-namespace
  labels:
    team: data-engineering
    env: test
  annotations:
    description: "Test namespace for data engineering team"
`

	var namespace Namespace
	err := yaml.Unmarshal([]byte(yamlContent), &namespace)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Test metadata
	if namespace.Metadata.Name != "test-namespace" {
		t.Errorf("Expected name 'test-namespace', got '%s'", namespace.Metadata.Name)
	}
	if namespace.Metadata.Labels["team"] != "data-engineering" {
		t.Errorf("Expected label 'team' to be 'data-engineering'")
	}
	if namespace.Metadata.Annotations["description"] != "Test namespace for data engineering team" {
		t.Errorf("Expected annotation description")
	}

	t.Log("Namespace parsing test passed!")
}

func TestNamespaceListParsing(t *testing.T) {
	yamlContent := `
items:
  - metadata:
      name: namespace1
      labels:
        env: dev
    spec:
      roleBindings:
        - role: owner
          members:
            - user:owner1@example.com
  - metadata:
      name: namespace2
      labels:
        env: prod
    spec:
      roleBindings:
        - role: owner
          members:
            - user:owner2@example.com
`

	var namespaceList NamespaceList
	err := yaml.Unmarshal([]byte(yamlContent), &namespaceList)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(namespaceList.Items) != 2 {
		t.Fatalf("Expected 2 namespaces, got %d", len(namespaceList.Items))
	}

	// Test first namespace
	if namespaceList.Items[0].Metadata.Name != "namespace1" {
		t.Errorf("Expected first namespace name 'namespace1', got '%s'", namespaceList.Items[0].Metadata.Name)
	}
	if namespaceList.Items[0].Metadata.Labels["env"] != "dev" {
		t.Errorf("Expected first namespace env label 'dev'")
	}

	// Test second namespace
	if namespaceList.Items[1].Metadata.Name != "namespace2" {
		t.Errorf("Expected second namespace name 'namespace2', got '%s'", namespaceList.Items[1].Metadata.Name)
	}
	if namespaceList.Items[1].Metadata.Labels["env"] != "prod" {
		t.Errorf("Expected second namespace env label 'prod'")
	}

	t.Log("NamespaceList parsing test passed!")
}

func TestNamespaceMinimal(t *testing.T) {
	yamlContent := `
metadata:
  name: minimal-namespace
`

	var namespace Namespace
	err := yaml.Unmarshal([]byte(yamlContent), &namespace)
	if err != nil {
		t.Fatalf("Failed to unmarshal minimal YAML: %v", err)
	}

	if namespace.Metadata.Name != "minimal-namespace" {
		t.Errorf("Expected name 'minimal-namespace', got '%s'", namespace.Metadata.Name)
	}

	// Should handle empty/nil spec gracefully
	if namespace.Spec.RoleBindings != nil && len(namespace.Spec.RoleBindings) > 0 {
		t.Errorf("Expected no role bindings in minimal spec")
	}

	t.Log("Minimal Namespace parsing test passed!")
}
