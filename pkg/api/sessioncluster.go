package api

import (
	"fmt"
	"time"
)

// SessionCluster represents a VVP session cluster
type SessionCluster struct {
	APIVersion string                 `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string                 `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   SessionClusterMetadata `json:"metadata" yaml:"metadata"`
	Spec       SessionClusterSpec     `json:"spec" yaml:"spec"`
	Status     SessionClusterStatus   `json:"status,omitempty" yaml:"status,omitempty"`
}

// SessionClusterMetadata holds session cluster metadata
type SessionClusterMetadata struct {
	ID              string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name            string            `json:"name" yaml:"name"`
	Namespace       string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty" yaml:"createdAt,omitempty"`
	ModifiedAt      time.Time         `json:"modifiedAt,omitempty" yaml:"modifiedAt,omitempty"`
	ResourceVersion int32             `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
}

// SessionClusterSpec holds session cluster specification
type SessionClusterSpec struct {
	DeploymentTargetName string                  `json:"deploymentTargetName" yaml:"deploymentTargetName"`
	FlinkConfiguration   map[string]string       `json:"flinkConfiguration,omitempty" yaml:"flinkConfiguration,omitempty"`
	FlinkImageRegistry   string                  `json:"flinkImageRegistry" yaml:"flinkImageRegistry"`
	FlinkImageRepository string                  `json:"flinkImageRepository" yaml:"flinkImageRepository"`
	FlinkImageTag        string                  `json:"flinkImageTag" yaml:"flinkImageTag"`
	FlinkVersion         string                  `json:"flinkVersion" yaml:"flinkVersion"`
	Kubernetes           interface{}             `json:"kubernetes,omitempty" yaml:"kubernetes,omitempty"`
	Logging              interface{}             `json:"logging,omitempty" yaml:"logging,omitempty"`
	NumberOfTaskManagers int32                   `json:"numberOfTaskManagers" yaml:"numberOfTaskManagers"`
	Resources            map[string]ResourceSpec `json:"resources" yaml:"resources"`
	State                string                  `json:"state" yaml:"state"` // STOPPED, RUNNING
}

// SessionClusterStatus holds session cluster status
type SessionClusterStatus struct {
	State   string                       `json:"state,omitempty" yaml:"state,omitempty"` // STOPPED, STARTING, RUNNING, UPDATING, STOPPING, FAILED
	Running *SessionClusterStatusRunning `json:"running,omitempty" yaml:"running,omitempty"`
	Failure *Failure                     `json:"failure,omitempty" yaml:"failure,omitempty"`
}

// SessionClusterStatusRunning holds running session cluster status
type SessionClusterStatusRunning struct {
	TransitionTime time.Time `json:"transitionTime,omitempty" yaml:"transitionTime,omitempty"`
}

// Failure holds failure information
type Failure struct {
	Message string    `json:"message,omitempty" yaml:"message,omitempty"`
	Reason  string    `json:"reason,omitempty" yaml:"reason,omitempty"`
	Time    time.Time `json:"time,omitempty" yaml:"time,omitempty"`
}

// SessionClusterList represents a list of session clusters
type SessionClusterList struct {
	APIVersion string           `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string           `json:"kind,omitempty" yaml:"kind,omitempty"`
	Items      []SessionCluster `json:"items" yaml:"items"`
}

// ListSessionClusters lists all session clusters in a namespace
func (c *Client) ListSessionClusters(namespace string) (*SessionClusterList, error) {
	var result SessionClusterList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/sessionclusters", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSessionCluster gets a session cluster by name
func (c *Client) GetSessionCluster(namespace, name string) (*SessionCluster, error) {
	var result SessionCluster
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/sessionclusters/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateSessionCluster creates a new session cluster
func (c *Client) CreateSessionCluster(namespace string, sessionCluster *SessionCluster) (*SessionCluster, error) {
	var result SessionCluster
	resp, err := c.httpClient.R().
		SetBody(sessionCluster).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/namespaces/%s/sessionclusters", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateSessionCluster updates an existing session cluster (PATCH)
func (c *Client) UpdateSessionCluster(namespace, name string, sessionCluster *SessionCluster) (*SessionCluster, error) {
	var result SessionCluster
	resp, err := c.httpClient.R().
		SetBody(sessionCluster).
		SetResult(&result).
		Patch(fmt.Sprintf("/api/v1/namespaces/%s/sessionclusters/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpsertSessionCluster creates or replaces a session cluster (PUT)
func (c *Client) UpsertSessionCluster(namespace, name string, sessionCluster *SessionCluster) (*SessionCluster, error) {
	var result SessionCluster
	resp, err := c.httpClient.R().
		SetBody(sessionCluster).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/namespaces/%s/sessionclusters/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteSessionCluster deletes a session cluster
func (c *Client) DeleteSessionCluster(namespace, name string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/v1/namespaces/%s/sessionclusters/%s", namespace, name))

	return handleResponse(resp, err)
}
