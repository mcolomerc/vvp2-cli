package api

import (
	"fmt"
	"time"
)

// Savepoint represents a VVP savepoint
type Savepoint struct {
	APIVersion string            `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string            `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   SavepointMetadata `json:"metadata" yaml:"metadata"`
	Spec       SavepointSpec     `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     SavepointStatus   `json:"status,omitempty" yaml:"status,omitempty"`
}

// SavepointMetadata holds savepoint metadata
type SavepointMetadata struct {
	ID              string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace       string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty" yaml:"createdAt,omitempty"`
	ModifiedAt      time.Time         `json:"modifiedAt,omitempty" yaml:"modifiedAt,omitempty"`
	ResourceVersion int32             `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
}

// SavepointSpec holds savepoint specification
type SavepointSpec struct {
	DeploymentID string `json:"deploymentId,omitempty" yaml:"deploymentId,omitempty"`
	JobID        string `json:"jobId,omitempty" yaml:"jobId,omitempty"`
}

// SavepointStatus holds savepoint status information
type SavepointStatus struct {
	State     string                  `json:"state,omitempty" yaml:"state,omitempty"`
	Completed *SavepointStatusDetails `json:"completed,omitempty" yaml:"completed,omitempty"`
	Failed    *SavepointStatusDetails `json:"failed,omitempty" yaml:"failed,omitempty"`
}

// SavepointStatusDetails holds savepoint status details
type SavepointStatusDetails struct {
	Location string    `json:"location,omitempty" yaml:"location,omitempty"`
	Time     time.Time `json:"time,omitempty" yaml:"time,omitempty"`
	Message  string    `json:"message,omitempty" yaml:"message,omitempty"`
	Reason   string    `json:"reason,omitempty" yaml:"reason,omitempty"`
}

// SavepointList represents a list of savepoints
type SavepointList struct {
	APIVersion string      `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string      `json:"kind,omitempty" yaml:"kind,omitempty"`
	Items      []Savepoint `json:"items" yaml:"items"`
}

// SavepointCreationRequest represents a request to create a savepoint
type SavepointCreationRequest struct {
	Metadata SavepointMetadata `json:"metadata" yaml:"metadata"`
	Spec     SavepointSpec     `json:"spec" yaml:"spec"`
}

// ListSavepoints lists all savepoints in a namespace
func (c *Client) ListSavepoints(namespace string) (*SavepointList, error) {
	var result SavepointList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/savepoints", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSavepoint gets a savepoint by ID
func (c *Client) GetSavepoint(namespace, savepointID string) (*Savepoint, error) {
	var result Savepoint
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/savepoints/%s", namespace, savepointID))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateSavepoint creates a new savepoint
func (c *Client) CreateSavepoint(namespace string, savepoint *SavepointCreationRequest) (*Savepoint, error) {
	var result Savepoint
	resp, err := c.httpClient.R().
		SetBody(savepoint).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/namespaces/%s/savepoints", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteSavepoint deletes a savepoint
func (c *Client) DeleteSavepoint(namespace, savepointID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/v1/namespaces/%s/savepoints/%s", namespace, savepointID))

	return handleResponse(resp, err)
}
