package api

import (
	"fmt"
)

// Status represents the VVP platform status
type Status struct {
	APIVersion    string         `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind          string         `json:"kind,omitempty" yaml:"kind,omitempty"`
	Health        HealthStatus   `json:"health,omitempty" yaml:"health,omitempty"`
	Version       VersionInfo    `json:"version,omitempty" yaml:"version,omitempty"`
	Components    []Component    `json:"components,omitempty" yaml:"components,omitempty"`
	ResourceUsage *ResourceUsage `json:"resourceUsage,omitempty" yaml:"resourceUsage,omitempty"`
}

// HealthStatus represents the health status of the platform
type HealthStatus struct {
	Status  string `json:"status,omitempty" yaml:"status,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

// VersionInfo represents version information
type VersionInfo struct {
	Platform   string `json:"platform,omitempty" yaml:"platform,omitempty"`
	Flink      string `json:"flink,omitempty" yaml:"flink,omitempty"`
	BuildTime  string `json:"buildTime,omitempty" yaml:"buildTime,omitempty"`
	CommitHash string `json:"commitHash,omitempty" yaml:"commitHash,omitempty"`
	Edition    string `json:"edition,omitempty" yaml:"edition,omitempty"`
}

// Component represents a platform component status
type Component struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Status  string `json:"status,omitempty" yaml:"status,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// ResourceUsage represents resource usage information
type ResourceUsage struct {
	Deployments     int `json:"deployments,omitempty" yaml:"deployments,omitempty"`
	Jobs            int `json:"jobs,omitempty" yaml:"jobs,omitempty"`
	SessionClusters int `json:"sessionClusters,omitempty" yaml:"sessionClusters,omitempty"`
	Namespaces      int `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`
}

// GetStatus retrieves the platform status
func (c *Client) GetStatus() (*Status, error) {
	var result Status
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/api/v1/status")

	if err := handleResponse(resp, err); err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	return &result, nil
}
