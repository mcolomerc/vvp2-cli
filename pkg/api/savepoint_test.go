package api

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSavepointParsing(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: Savepoint
metadata:
  name: savepoint-20231128-120000
  namespace: default
  labels:
    deployment: my-flink-job
    trigger: manual
spec:
  deploymentId: deployment-123
  jobId: job-456
status:
  state: COMPLETED
  completed:
    location: "s3://my-bucket/savepoints/savepoint-20231128-120000"
    time: "2023-11-28T12:00:00Z"
`

	var savepoint Savepoint
	err := yaml.Unmarshal([]byte(yamlContent), &savepoint)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Test metadata
	if savepoint.Metadata.Name != "savepoint-20231128-120000" {
		t.Errorf("Expected name 'savepoint-20231128-120000', got '%s'", savepoint.Metadata.Name)
	}
	if savepoint.Metadata.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got '%s'", savepoint.Metadata.Namespace)
	}
	if savepoint.Metadata.Labels["deployment"] != "my-flink-job" {
		t.Errorf("Expected label 'deployment' to be 'my-flink-job'")
	}

	// Test spec
	if savepoint.Spec.DeploymentID != "deployment-123" {
		t.Errorf("Expected deploymentId 'deployment-123', got '%s'", savepoint.Spec.DeploymentID)
	}
	if savepoint.Spec.JobID != "job-456" {
		t.Errorf("Expected jobId 'job-456', got '%s'", savepoint.Spec.JobID)
	}

	// Test status
	if savepoint.Status.State != "COMPLETED" {
		t.Errorf("Expected state 'COMPLETED', got '%s'", savepoint.Status.State)
	}
	if savepoint.Status.Completed == nil {
		t.Fatal("Expected completed status to be non-nil")
	}
	if savepoint.Status.Completed.Location != "s3://my-bucket/savepoints/savepoint-20231128-120000" {
		t.Errorf("Expected location, got '%s'", savepoint.Status.Completed.Location)
	}

	t.Log("Savepoint parsing test passed!")
}

func TestSavepointListParsing(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: SavepointList
items:
  - apiVersion: v1
    kind: Savepoint
    metadata:
      name: savepoint-1
      namespace: default
    spec:
      deploymentId: deployment-1
      jobId: job-1
    status:
      state: COMPLETED
  - apiVersion: v1
    kind: Savepoint
    metadata:
      name: savepoint-2
      namespace: default
    spec:
      deploymentId: deployment-2
      jobId: job-2
    status:
      state: IN_PROGRESS
`

	var savepointList SavepointList
	err := yaml.Unmarshal([]byte(yamlContent), &savepointList)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(savepointList.Items) != 2 {
		t.Fatalf("Expected 2 savepoints, got %d", len(savepointList.Items))
	}

	// Test first savepoint
	if savepointList.Items[0].Metadata.Name != "savepoint-1" {
		t.Errorf("Expected first savepoint name 'savepoint-1', got '%s'", savepointList.Items[0].Metadata.Name)
	}
	if savepointList.Items[0].Status.State != "COMPLETED" {
		t.Errorf("Expected first savepoint state 'COMPLETED', got '%s'", savepointList.Items[0].Status.State)
	}

	// Test second savepoint
	if savepointList.Items[1].Metadata.Name != "savepoint-2" {
		t.Errorf("Expected second savepoint name 'savepoint-2', got '%s'", savepointList.Items[1].Metadata.Name)
	}
	if savepointList.Items[1].Status.State != "IN_PROGRESS" {
		t.Errorf("Expected second savepoint state 'IN_PROGRESS', got '%s'", savepointList.Items[1].Status.State)
	}

	t.Log("SavepointList parsing test passed!")
}

func TestSavepointCreationRequest(t *testing.T) {
	yamlContent := `
metadata:
  name: manual-savepoint
  namespace: default
  labels:
    trigger: manual
spec:
  deploymentId: my-deployment
  jobId: my-job-id
`

	var request SavepointCreationRequest
	err := yaml.Unmarshal([]byte(yamlContent), &request)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if request.Metadata.Name != "manual-savepoint" {
		t.Errorf("Expected name 'manual-savepoint', got '%s'", request.Metadata.Name)
	}
	if request.Spec.DeploymentID != "my-deployment" {
		t.Errorf("Expected deploymentId 'my-deployment', got '%s'", request.Spec.DeploymentID)
	}

	t.Log("SavepointCreationRequest parsing test passed!")
}

func TestSavepointFailedStatus(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: Savepoint
metadata:
  name: failed-savepoint
  namespace: default
spec:
  deploymentId: deployment-123
  jobId: job-456
status:
  state: FAILED
  failed:
    message: "Failed to create savepoint"
    reason: "TimeoutException"
    time: "2023-11-28T12:00:00Z"
`

	var savepoint Savepoint
	err := yaml.Unmarshal([]byte(yamlContent), &savepoint)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if savepoint.Status.State != "FAILED" {
		t.Errorf("Expected state 'FAILED', got '%s'", savepoint.Status.State)
	}
	if savepoint.Status.Failed == nil {
		t.Fatal("Expected failed status to be non-nil")
	}
	if savepoint.Status.Failed.Message != "Failed to create savepoint" {
		t.Errorf("Expected failure message, got '%s'", savepoint.Status.Failed.Message)
	}
	if savepoint.Status.Failed.Reason != "TimeoutException" {
		t.Errorf("Expected failure reason 'TimeoutException', got '%s'", savepoint.Status.Failed.Reason)
	}

	t.Log("Failed Savepoint parsing test passed!")
}
