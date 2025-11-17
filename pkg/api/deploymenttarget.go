package api

import (
	"fmt"
	"time"
)

// DeploymentTargetResource represents a VVP deployment target
type DeploymentTargetResource struct {
	Metadata DeploymentTargetMetadata `json:"metadata"`
	Spec     DeploymentTargetSpec     `json:"spec"`
	Status   DeploymentTargetStatus   `json:"status,omitempty"`
}

// DeploymentTargetMetadata holds deployment target metadata
type DeploymentTargetMetadata struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"createdAt,omitempty"`
	ModifiedAt  time.Time         `json:"modifiedAt,omitempty"`
}

// DeploymentTargetSpec holds deployment target specification
type DeploymentTargetSpec struct {
	Kubernetes KubernetesTarget `json:"kubernetes,omitempty"`
}

// KubernetesTarget defines Kubernetes-specific settings
type KubernetesTarget struct {
	Namespace string `json:"namespace,omitempty"`
}

// DeploymentTargetStatus holds deployment target status
type DeploymentTargetStatus struct {
	State string `json:"state,omitempty"`
}

// DeploymentTargetList represents a list of deployment targets
type DeploymentTargetList struct {
	Items []DeploymentTargetResource `json:"items"`
}

// ListDeploymentTargets lists all deployment targets in a namespace
func (c *Client) ListDeploymentTargets(namespace string) (*DeploymentTargetList, error) {
	var result DeploymentTargetList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/deployment-targets", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDeploymentTarget gets a deployment target by name
func (c *Client) GetDeploymentTarget(namespace, name string) (*DeploymentTargetResource, error) {
	var result DeploymentTargetResource
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/deployment-targets/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateDeploymentTarget creates a new deployment target
func (c *Client) CreateDeploymentTarget(namespace string, target *DeploymentTargetResource) (*DeploymentTargetResource, error) {
	var result DeploymentTargetResource
	resp, err := c.httpClient.R().
		SetBody(target).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/namespaces/%s/deployment-targets", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateDeploymentTarget updates an existing deployment target
func (c *Client) UpdateDeploymentTarget(namespace, name string, target *DeploymentTargetResource) (*DeploymentTargetResource, error) {
	var result DeploymentTargetResource
	resp, err := c.httpClient.R().
		SetBody(target).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/namespaces/%s/deployment-targets/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteDeploymentTarget deletes a deployment target
func (c *Client) DeleteDeploymentTarget(namespace, name string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/v1/namespaces/%s/deployment-targets/%s", namespace, name))

	return handleResponse(resp, err)
}
