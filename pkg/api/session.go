package api

import (
	"fmt"
	"time"
)

// NOTE: Session endpoints do not exist in VVP API.
// The endpoint /api/v1/namespaces/{ns}/sessions returns 404/500 errors.
// Sessions are called "session clusters" in VVP and are referenced in deployments
// via the sessionClusterName field, but there's no standalone API for managing them.
// These functions are kept for reference but will not work.

// Session represents a VVP session (SQL session)
type Session struct {
	Metadata SessionMetadata `json:"metadata" yaml:"metadata"`
	Spec     SessionSpec     `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status   SessionStatus   `json:"status,omitempty" yaml:"status,omitempty"`
}

// SessionMetadata holds session metadata
type SessionMetadata struct {
	ID          string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace" yaml:"namespace"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"createdAt,omitempty" yaml:"createdAt,omitempty"`
	ModifiedAt  time.Time         `json:"modifiedAt,omitempty" yaml:"modifiedAt,omitempty"`
}

// SessionSpec holds session specification
type SessionSpec struct {
	DeploymentTargetID            string                        `json:"deploymentTargetId,omitempty" yaml:"deploymentTargetId,omitempty"`
	FlinkConfiguration            map[string]string             `json:"flinkConfiguration,omitempty" yaml:"flinkConfiguration,omitempty"`
	FlinkVersion                  string                        `json:"flinkVersion,omitempty" yaml:"flinkVersion,omitempty"`
	SessionClusterResourceProfile SessionClusterResourceProfile `json:"sessionClusterResourceProfile,omitempty" yaml:"sessionClusterResourceProfile,omitempty"`
}

// SessionClusterResourceProfile defines resource profile for session cluster
type SessionClusterResourceProfile struct {
	CPU    string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
}

// SessionStatus holds session status
type SessionStatus struct {
	State string `json:"state,omitempty" yaml:"state,omitempty"`
}

// SessionList represents a list of sessions
type SessionList struct {
	Items []Session `json:"items"`
}

// ListSessions lists all sessions in a namespace
// WARNING: This endpoint does not exist in VVP and will return 404 errors
func (c *Client) ListSessions(namespace string) (*SessionList, error) {
	var result SessionList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/sessions", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSession gets a session by name
// WARNING: This endpoint does not exist in VVP and will return 404 errors
func (c *Client) GetSession(namespace, name string) (*Session, error) {
	var result Session
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/sessions/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateSession creates a new session
// WARNING: This endpoint does not exist in VVP and will return 500 errors
func (c *Client) CreateSession(namespace string, session *Session) (*Session, error) {
	var result Session
	resp, err := c.httpClient.R().
		SetBody(session).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/namespaces/%s/sessions", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateSession updates an existing session
// WARNING: This endpoint does not exist in VVP and will return 500 errors
func (c *Client) UpdateSession(namespace, name string, session *Session) (*Session, error) {
	var result Session
	resp, err := c.httpClient.R().
		SetBody(session).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/namespaces/%s/sessions/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteSession deletes a session
// WARNING: This endpoint does not exist in VVP and will return 404 errors
func (c *Client) DeleteSession(namespace, name string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/v1/namespaces/%s/sessions/%s", namespace, name))

	return handleResponse(resp, err)
}
