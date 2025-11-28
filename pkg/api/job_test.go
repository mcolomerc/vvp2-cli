package api

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestJobParsing(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: Job
metadata:
  id: job-123
  name: my-flink-job
  namespace: default
  labels:
    app: flink
    deployment: my-deployment
spec:
  deploymentId: deployment-456
  state: RUNNING
status:
  state: RUNNING
  running:
    startTime: "2023-11-28T10:00:00Z"
    transitionTime: "2023-11-28T10:01:00Z"
    jobId: "flink-job-789"
`

	var job Job
	err := yaml.Unmarshal([]byte(yamlContent), &job)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Test metadata
	if job.Metadata.ID != "job-123" {
		t.Errorf("Expected id 'job-123', got '%s'", job.Metadata.ID)
	}
	if job.Metadata.Name != "my-flink-job" {
		t.Errorf("Expected name 'my-flink-job', got '%s'", job.Metadata.Name)
	}
	if job.Metadata.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got '%s'", job.Metadata.Namespace)
	}
	if job.Metadata.Labels["app"] != "flink" {
		t.Errorf("Expected label 'app' to be 'flink'")
	}

	// Test spec
	if job.Spec.DeploymentID != "deployment-456" {
		t.Errorf("Expected deploymentId 'deployment-456', got '%s'", job.Spec.DeploymentID)
	}
	if job.Spec.State != "RUNNING" {
		t.Errorf("Expected state 'RUNNING', got '%s'", job.Spec.State)
	}

	// Test status
	if job.Status.State != "RUNNING" {
		t.Errorf("Expected status state 'RUNNING', got '%s'", job.Status.State)
	}
	if job.Status.Running == nil {
		t.Fatal("Expected running status to be non-nil")
	}
	if job.Status.Running.JobID != "flink-job-789" {
		t.Errorf("Expected jobId 'flink-job-789', got '%s'", job.Status.Running.JobID)
	}

	t.Log("Job parsing test passed!")
}

func TestJobListParsing(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: JobList
items:
  - apiVersion: v1
    kind: Job
    metadata:
      id: job-1
      name: job-one
      namespace: default
    spec:
      deploymentId: deployment-1
      state: RUNNING
    status:
      state: RUNNING
  - apiVersion: v1
    kind: Job
    metadata:
      id: job-2
      name: job-two
      namespace: default
    spec:
      deploymentId: deployment-2
      state: CANCELLED
    status:
      state: CANCELLED
`

	var jobList JobList
	err := yaml.Unmarshal([]byte(yamlContent), &jobList)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(jobList.Items) != 2 {
		t.Fatalf("Expected 2 jobs, got %d", len(jobList.Items))
	}

	// Test first job
	if jobList.Items[0].Metadata.Name != "job-one" {
		t.Errorf("Expected first job name 'job-one', got '%s'", jobList.Items[0].Metadata.Name)
	}
	if jobList.Items[0].Status.State != "RUNNING" {
		t.Errorf("Expected first job state 'RUNNING', got '%s'", jobList.Items[0].Status.State)
	}

	// Test second job
	if jobList.Items[1].Metadata.Name != "job-two" {
		t.Errorf("Expected second job name 'job-two', got '%s'", jobList.Items[1].Metadata.Name)
	}
	if jobList.Items[1].Status.State != "CANCELLED" {
		t.Errorf("Expected second job state 'CANCELLED', got '%s'", jobList.Items[1].Status.State)
	}

	t.Log("JobList parsing test passed!")
}

func TestJobFailedStatus(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: Job
metadata:
  id: failed-job
  name: failed-job
  namespace: default
spec:
  deploymentId: deployment-123
  state: FAILED
status:
  state: FAILED
  failed:
    failureTime: "2023-11-28T12:00:00Z"
    reason: "Exception"
    message: "Job failed due to exception in user code"
`

	var job Job
	err := yaml.Unmarshal([]byte(yamlContent), &job)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if job.Status.State != "FAILED" {
		t.Errorf("Expected state 'FAILED', got '%s'", job.Status.State)
	}
	if job.Status.Failed == nil {
		t.Fatal("Expected failed status to be non-nil")
	}
	if job.Status.Failed.Reason != "Exception" {
		t.Errorf("Expected failure reason 'Exception', got '%s'", job.Status.Failed.Reason)
	}
	if job.Status.Failed.Message != "Job failed due to exception in user code" {
		t.Errorf("Expected failure message, got '%s'", job.Status.Failed.Message)
	}

	t.Log("Failed Job parsing test passed!")
}

func TestJobSuspendedStatus(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: Job
metadata:
  id: suspended-job
  name: suspended-job
  namespace: default
spec:
  deploymentId: deployment-123
  state: SUSPENDED
status:
  state: SUSPENDED
  suspended:
    suspensionTime: "2023-11-28T12:00:00Z"
`

	var job Job
	err := yaml.Unmarshal([]byte(yamlContent), &job)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if job.Status.State != "SUSPENDED" {
		t.Errorf("Expected state 'SUSPENDED', got '%s'", job.Status.State)
	}
	if job.Status.Suspended == nil {
		t.Fatal("Expected suspended status to be non-nil")
	}

	t.Log("Suspended Job parsing test passed!")
}

func TestJobFinishedStatus(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: Job
metadata:
  id: finished-job
  name: finished-job
  namespace: default
spec:
  deploymentId: deployment-123
  state: FINISHED
status:
  state: FINISHED
  finished:
    completionTime: "2023-11-28T12:00:00Z"
`

	var job Job
	err := yaml.Unmarshal([]byte(yamlContent), &job)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if job.Status.State != "FINISHED" {
		t.Errorf("Expected state 'FINISHED', got '%s'", job.Status.State)
	}
	if job.Status.Finished == nil {
		t.Fatal("Expected finished status to be non-nil")
	}

	t.Log("Finished Job parsing test passed!")
}
