package api

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSecretValueParsing(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: SecretValue
metadata:
  name: kafka-password
  namespace: default
  labels:
    app: kafka
    env: production
  annotations:
    description: "Kafka cluster password"
spec:
  kind: GENERIC
  value: "my-secret-password"
`

	var secretValue SecretValue
	err := yaml.Unmarshal([]byte(yamlContent), &secretValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Test metadata
	if secretValue.Metadata.Name != "kafka-password" {
		t.Errorf("Expected name 'kafka-password', got '%s'", secretValue.Metadata.Name)
	}
	if secretValue.Metadata.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got '%s'", secretValue.Metadata.Namespace)
	}
	if secretValue.Metadata.Labels["app"] != "kafka" {
		t.Errorf("Expected label 'app' to be 'kafka'")
	}
	if secretValue.Metadata.Annotations["description"] != "Kafka cluster password" {
		t.Errorf("Expected annotation description")
	}

	// Test spec
	if secretValue.Spec.Kind != "GENERIC" {
		t.Errorf("Expected kind 'GENERIC', got '%s'", secretValue.Spec.Kind)
	}
	if secretValue.Spec.Value != "my-secret-password" {
		t.Errorf("Expected value 'my-secret-password', got '%s'", secretValue.Spec.Value)
	}

	t.Log("SecretValue parsing test passed!")
}

func TestSecretValueListParsing(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: SecretValueList
items:
  - apiVersion: v1
    kind: SecretValue
    metadata:
      name: db-password
      namespace: default
    spec:
      kind: GENERIC
      value: "db-secret-123"
  - apiVersion: v1
    kind: SecretValue
    metadata:
      name: api-key
      namespace: default
    spec:
      kind: GENERIC
      value: "api-key-456"
`

	var secretValueList SecretValueList
	err := yaml.Unmarshal([]byte(yamlContent), &secretValueList)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(secretValueList.Items) != 2 {
		t.Fatalf("Expected 2 secret values, got %d", len(secretValueList.Items))
	}

	// Test first secret value
	if secretValueList.Items[0].Metadata.Name != "db-password" {
		t.Errorf("Expected first secret name 'db-password', got '%s'", secretValueList.Items[0].Metadata.Name)
	}
	if secretValueList.Items[0].Spec.Value != "db-secret-123" {
		t.Errorf("Expected first secret value 'db-secret-123', got '%s'", secretValueList.Items[0].Spec.Value)
	}

	// Test second secret value
	if secretValueList.Items[1].Metadata.Name != "api-key" {
		t.Errorf("Expected second secret name 'api-key', got '%s'", secretValueList.Items[1].Metadata.Name)
	}
	if secretValueList.Items[1].Spec.Value != "api-key-456" {
		t.Errorf("Expected second secret value 'api-key-456', got '%s'", secretValueList.Items[1].Spec.Value)
	}

	t.Log("SecretValueList parsing test passed!")
}

func TestSecretValueMinimal(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: SecretValue
metadata:
  name: minimal-secret
  namespace: default
spec:
  kind: GENERIC
  value: "secret-value"
`

	var secretValue SecretValue
	err := yaml.Unmarshal([]byte(yamlContent), &secretValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal minimal YAML: %v", err)
	}

	if secretValue.Metadata.Name != "minimal-secret" {
		t.Errorf("Expected name 'minimal-secret', got '%s'", secretValue.Metadata.Name)
	}
	if secretValue.Spec.Kind != "GENERIC" {
		t.Errorf("Expected kind 'GENERIC', got '%s'", secretValue.Spec.Kind)
	}
	if secretValue.Spec.Value != "secret-value" {
		t.Errorf("Expected value 'secret-value', got '%s'", secretValue.Spec.Value)
	}

	t.Log("Minimal SecretValue parsing test passed!")
}

func TestSecretValueMultipleLabels(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: SecretValue
metadata:
  name: labeled-secret
  namespace: default
  labels:
    team: data-engineering
    env: production
    app: flink
    managed-by: vvp2-cli
spec:
  kind: GENERIC
  value: "complex-secret"
`

	var secretValue SecretValue
	err := yaml.Unmarshal([]byte(yamlContent), &secretValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(secretValue.Metadata.Labels) != 4 {
		t.Errorf("Expected 4 labels, got %d", len(secretValue.Metadata.Labels))
	}

	expectedLabels := map[string]string{
		"team":       "data-engineering",
		"env":        "production",
		"app":        "flink",
		"managed-by": "vvp2-cli",
	}

	for key, expectedValue := range expectedLabels {
		if value, ok := secretValue.Metadata.Labels[key]; !ok {
			t.Errorf("Expected label '%s' to exist", key)
		} else if value != expectedValue {
			t.Errorf("Expected label '%s' to be '%s', got '%s'", key, expectedValue, value)
		}
	}

	t.Log("SecretValue with multiple labels parsing test passed!")
}
