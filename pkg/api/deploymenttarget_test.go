package api

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDeploymentTargetParsing(t *testing.T) {
	yamlContent := `
metadata:
  name: kubernetes-target
  namespace: default
  labels:
    env: production
    region: us-west-2
  annotations:
    description: "Production Kubernetes target"
spec:
  kubernetes:
    namespace: vvp-jobs
`

	var target DeploymentTargetResource
	err := yaml.Unmarshal([]byte(yamlContent), &target)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Test metadata
	if target.Metadata.Name != "kubernetes-target" {
		t.Errorf("Expected name 'kubernetes-target', got '%s'", target.Metadata.Name)
	}
	if target.Metadata.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got '%s'", target.Metadata.Namespace)
	}
	if target.Metadata.Labels["env"] != "production" {
		t.Errorf("Expected label 'env' to be 'production'")
	}
	if target.Metadata.Annotations["description"] != "Production Kubernetes target" {
		t.Errorf("Expected annotation description")
	}

	// Test spec
	if target.Spec.Kubernetes.Namespace != "vvp-jobs" {
		t.Errorf("Expected kubernetes namespace 'vvp-jobs', got '%s'", target.Spec.Kubernetes.Namespace)
	}

	t.Log("DeploymentTarget parsing test passed!")
}

func TestDeploymentTargetListParsing(t *testing.T) {
	yamlContent := `
items:
  - metadata:
      name: dev-target
      namespace: default
    spec:
      kubernetes:
        namespace: vvp-dev
  - metadata:
      name: prod-target
      namespace: default
    spec:
      kubernetes:
        namespace: vvp-prod
`

	var targetList DeploymentTargetList
	err := yaml.Unmarshal([]byte(yamlContent), &targetList)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(targetList.Items) != 2 {
		t.Fatalf("Expected 2 deployment targets, got %d", len(targetList.Items))
	}

	// Test first target
	if targetList.Items[0].Metadata.Name != "dev-target" {
		t.Errorf("Expected first target name 'dev-target', got '%s'", targetList.Items[0].Metadata.Name)
	}
	if targetList.Items[0].Spec.Kubernetes.Namespace != "vvp-dev" {
		t.Errorf("Expected first target k8s namespace 'vvp-dev', got '%s'", targetList.Items[0].Spec.Kubernetes.Namespace)
	}

	// Test second target
	if targetList.Items[1].Metadata.Name != "prod-target" {
		t.Errorf("Expected second target name 'prod-target', got '%s'", targetList.Items[1].Metadata.Name)
	}
	if targetList.Items[1].Spec.Kubernetes.Namespace != "vvp-prod" {
		t.Errorf("Expected second target k8s namespace 'vvp-prod', got '%s'", targetList.Items[1].Spec.Kubernetes.Namespace)
	}

	t.Log("DeploymentTargetList parsing test passed!")
}

func TestDeploymentTargetMinimal(t *testing.T) {
	yamlContent := `
metadata:
  name: minimal-target
  namespace: default
spec:
  kubernetes:
    namespace: vvp-minimal
`

	var target DeploymentTargetResource
	err := yaml.Unmarshal([]byte(yamlContent), &target)
	if err != nil {
		t.Fatalf("Failed to unmarshal minimal YAML: %v", err)
	}

	if target.Metadata.Name != "minimal-target" {
		t.Errorf("Expected name 'minimal-target', got '%s'", target.Metadata.Name)
	}
	if target.Spec.Kubernetes.Namespace != "vvp-minimal" {
		t.Errorf("Expected k8s namespace 'vvp-minimal', got '%s'", target.Spec.Kubernetes.Namespace)
	}

	t.Log("Minimal DeploymentTarget parsing test passed!")
}
