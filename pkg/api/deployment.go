package api

import (
	"fmt"
	"time"
)

// Deployment represents a VVP deployment
type Deployment struct {
	Metadata DeploymentMetadata  `json:"metadata"`
	Spec     DeploymentSpec      `json:"spec"`
	Status   *DeploymentStatus   `json:"status,omitempty"`
}

// DeploymentMetadata holds deployment metadata
type DeploymentMetadata struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"createdAt,omitempty"`
	ModifiedAt  time.Time         `json:"modifiedAt,omitempty"`
}

// DeploymentSpec holds deployment specification
type DeploymentSpec struct {
	State              string           `json:"state"`
	UpgradeStrategy    UpgradeStrategy  `json:"upgradeStrategy,omitempty"`
	RestoreStrategy    RestoreStrategy  `json:"restoreStrategy,omitempty"`
	DeploymentTarget   DeploymentTarget `json:"deploymentTargetId,omitempty"`
	Template           Template         `json:"template"`
	MaxSavepointAge    string           `json:"maxSavepointCreationTime,omitempty"`
	MaxJobCreationTime string           `json:"maxJobCreationTime,omitempty"`
}

// DeploymentStatus holds deployment status
type DeploymentStatus struct {
	State   string `json:"state"`
	Running bool   `json:"running"`
}

// UpgradeStrategy defines upgrade strategy
type UpgradeStrategy struct {
	Kind string `json:"kind"`
}

// RestoreStrategy defines restore strategy
type RestoreStrategy struct {
	Kind                  string `json:"kind"`
	AllowNonRestoredState bool   `json:"allowNonRestoredState,omitempty"`
}

// DeploymentTarget defines where to deploy
type DeploymentTarget struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Template defines the Flink job template
type Template struct {
	Spec TemplateSpec `json:"spec"`
}

// TemplateSpec holds template specification
type TemplateSpec struct {
	Artifact             Artifact          `json:"artifact"`
	Parallelism          int               `json:"parallelism,omitempty"`
	NumberOfTaskManagers int               `json:"numberOfTaskManagers,omitempty"`
	Resources            Resources         `json:"resources,omitempty"`
	FlinkVersion         string            `json:"flinkVersion,omitempty"`
	FlinkImageTag        string            `json:"flinkImageTag,omitempty"`
	Logging              Logging           `json:"logging,omitempty"`
	FlinkConfiguration   map[string]string `json:"flinkConfiguration,omitempty"`
}

// Artifact represents a JAR artifact
type Artifact struct {
	Kind          string `json:"kind"`
	JarURI        string `json:"jarUri"`
	MainClass     string `json:"mainClass,omitempty"`
	EntryClass    string `json:"entryClass,omitempty"`
	MainArgs      string `json:"mainArgs,omitempty"`
	FlinkVersion  string `json:"flinkVersion,omitempty"`
	FlinkImageTag string `json:"flinkImageTag,omitempty"`
}

// Resources defines resource requirements
type Resources struct {
	JobManager  ResourceSpec `json:"jobmanager,omitempty"`
	TaskManager ResourceSpec `json:"taskmanager,omitempty"`
}

// ResourceSpec defines CPU and memory resources
// CPU and Memory can be strings or numbers from the API, so we use interface{}
type ResourceSpec struct {
	CPU    interface{} `json:"cpu,omitempty"`
	Memory interface{} `json:"memory,omitempty"`
}

// Logging defines logging configuration
type Logging struct {
	Log4jLoggers map[string]string `json:"log4jLoggers,omitempty"`
}

// DeploymentList represents a list of deployments
type DeploymentList struct {
	Items []DeploymentWithInfo `json:"items"`
}

// DeploymentWithInfo wraps deployment with operator info
type DeploymentWithInfo struct {
	Deployment Deployment `json:"deployment"`
}

// ListDeployments lists all deployments in a namespace
func (c *Client) ListDeployments(namespace string) (*DeploymentList, error) {
	var result DeploymentList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/deployments/with-cr", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDeployment gets a deployment by name
func (c *Client) GetDeployment(namespace, name string) (*Deployment, error) {
	var wrapper DeploymentWithInfo
	resp, err := c.httpClient.R().
		SetResult(&wrapper).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/deployments/with-cr/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &wrapper.Deployment, nil
}

// CreateDeployment creates a new deployment
func (c *Client) CreateDeployment(namespace string, deployment *Deployment) (*Deployment, error) {
	var result Deployment
	resp, err := c.httpClient.R().
		SetBody(deployment).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/namespaces/%s/deployments", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateDeployment updates an existing deployment
func (c *Client) UpdateDeployment(namespace, name string, deployment *Deployment) (*Deployment, error) {
	var result Deployment
	resp, err := c.httpClient.R().
		SetBody(deployment).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/namespaces/%s/deployments/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteDeployment deletes a deployment
func (c *Client) DeleteDeployment(namespace, name string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/v1/namespaces/%s/deployments/%s", namespace, name))

	return handleResponse(resp, err)
}

// UpdateDeploymentState updates the state of a deployment
func (c *Client) UpdateDeploymentState(namespace, name, state string) (*Deployment, error) {
	patch := map[string]interface{}{
		"spec": map[string]interface{}{
			"state": state,
		},
	}

	var result Deployment
	resp, err := c.httpClient.R().
		SetBody(patch).
		SetResult(&result).
		Patch(fmt.Sprintf("/api/v1/namespaces/%s/deployments/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}
