package api

import (
	"fmt"
	"time"
)

// Job represents a VVP job
type Job struct {
	APIVersion string      `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string      `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   JobMetadata `json:"metadata" yaml:"metadata"`
	Spec       JobSpec     `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     JobStatus   `json:"status,omitempty" yaml:"status,omitempty"`
}

// JobMetadata holds job metadata
type JobMetadata struct {
	ID              string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace       string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty" yaml:"createdAt,omitempty"`
	ModifiedAt      time.Time         `json:"modifiedAt,omitempty" yaml:"modifiedAt,omitempty"`
	ResourceVersion int32             `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
}

// JobSpec holds job specification
type JobSpec struct {
	DeploymentID string `json:"deploymentId,omitempty" yaml:"deploymentId,omitempty"`
	State        string `json:"state,omitempty" yaml:"state,omitempty"`
}

// JobStatus holds job status information
type JobStatus struct {
	State       string               `json:"state,omitempty" yaml:"state,omitempty"`
	Running     *JobStatusRunning    `json:"running,omitempty" yaml:"running,omitempty"`
	Failed      *JobStatusFailed     `json:"failed,omitempty" yaml:"failed,omitempty"`
	Cancelled   *JobStatusCancelled  `json:"cancelled,omitempty" yaml:"cancelled,omitempty"`
	Finished    *JobStatusFinished   `json:"finished,omitempty" yaml:"finished,omitempty"`
	Suspended   *JobStatusSuspended  `json:"suspended,omitempty" yaml:"suspended,omitempty"`
	Terminating *JobStatusTerminating `json:"terminating,omitempty" yaml:"terminating,omitempty"`
}

// JobStatusRunning holds running job status
type JobStatusRunning struct {
	StartTime      time.Time `json:"startTime,omitempty" yaml:"startTime,omitempty"`
	TransitionTime time.Time `json:"transitionTime,omitempty" yaml:"transitionTime,omitempty"`
	JobID          string    `json:"jobId,omitempty" yaml:"jobId,omitempty"`
}

// JobStatusFailed holds failed job status
type JobStatusFailed struct {
	FailureTime time.Time `json:"failureTime,omitempty" yaml:"failureTime,omitempty"`
	Reason      string    `json:"reason,omitempty" yaml:"reason,omitempty"`
	Message     string    `json:"message,omitempty" yaml:"message,omitempty"`
}

// JobStatusCancelled holds cancelled job status
type JobStatusCancelled struct {
	CancellationTime time.Time `json:"cancellationTime,omitempty" yaml:"cancellationTime,omitempty"`
}

// JobStatusFinished holds finished job status
type JobStatusFinished struct {
	CompletionTime time.Time `json:"completionTime,omitempty" yaml:"completionTime,omitempty"`
}

// JobStatusSuspended holds suspended job status
type JobStatusSuspended struct {
	SuspensionTime time.Time `json:"suspensionTime,omitempty" yaml:"suspensionTime,omitempty"`
}

// JobStatusTerminating holds terminating job status
type JobStatusTerminating struct {
	TransitionTime time.Time `json:"transitionTime,omitempty" yaml:"transitionTime,omitempty"`
}

// JobList represents a list of jobs
type JobList struct {
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Items      []Job  `json:"items" yaml:"items"`
}

// ListJobs lists all jobs in a namespace
func (c *Client) ListJobs(namespace string) (*JobList, error) {
	var result JobList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/jobs", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetJob gets a job by ID
func (c *Client) GetJob(namespace, jobID string) (*Job, error) {
	var result Job
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/jobs/%s", namespace, jobID))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}
