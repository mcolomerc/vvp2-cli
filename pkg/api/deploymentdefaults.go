package api

import (
	"fmt"
	"time"
)

// DeploymentDefaults represents the namespace-level defaults for deployments
type DeploymentDefaults struct {
	APIVersion string                     `json:"apiVersion,omitempty"`
	Kind       string                     `json:"kind,omitempty"`
	Metadata   DeploymentDefaultsMetadata `json:"metadata"`
	Spec       DeploymentSpec             `json:"spec"`
}

// DeploymentDefaultsMetadata holds metadata for deployment defaults
type DeploymentDefaultsMetadata struct {
	ID              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Namespace       string            `json:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty"`
	ModifiedAt      time.Time         `json:"modifiedAt,omitempty"`
	ResourceVersion int32             `json:"resourceVersion,omitempty"`
}

// SecretValue represents a secret value resource used by some endpoints
type SecretValue struct {
	APIVersion string              `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string              `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   SecretValueMetadata `json:"metadata" yaml:"metadata"`
	Spec       SecretValueSpec     `json:"spec" yaml:"spec"`
}

// SecretValueMetadata contains metadata for a secret value
type SecretValueMetadata struct {
	ID              string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name            string            `json:"name" yaml:"name"`
	Namespace       string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty" yaml:"createdAt,omitempty"`
	ModifiedAt      time.Time         `json:"modifiedAt,omitempty" yaml:"modifiedAt,omitempty"`
	ResourceVersion int32             `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
}

// SecretValueSpec contains the secret value specification
type SecretValueSpec struct {
	Kind  string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

// GetDeploymentDefaults retrieves the deployment defaults for a namespace
func (c *Client) GetDeploymentDefaults(namespace string) (*DeploymentDefaults, error) {
	var result DeploymentDefaults
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/deployment-defaults", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// ReplaceDeploymentDefaults replaces the deployment defaults for a namespace
func (c *Client) ReplaceDeploymentDefaults(namespace string, defaults *DeploymentDefaults) (*DeploymentDefaults, error) {
	var result DeploymentDefaults
	resp, err := c.httpClient.R().
		SetBody(defaults).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/namespaces/%s/deployment-defaults", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateDeploymentDefaults updates the deployment defaults via PATCH
// According to the Application Manager spec, this endpoint accepts a SecretValue body.
func (c *Client) UpdateDeploymentDefaults(namespace string, secret *SecretValue) (*DeploymentDefaults, error) {
	var result DeploymentDefaults
	resp, err := c.httpClient.R().
		SetBody(secret).
		SetResult(&result).
		Patch(fmt.Sprintf("/api/v1/namespaces/%s/deployment-defaults", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}
