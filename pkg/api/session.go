package api

import (
	"fmt"
	"time"
)

// Session represents a VVP session (SQL session)
type Session struct {
	Metadata SessionMetadata `json:"metadata"`
	Spec     SessionSpec     `json:"spec,omitempty"`
	Status   SessionStatus   `json:"status,omitempty"`
}

// SessionMetadata holds session metadata
type SessionMetadata struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"createdAt,omitempty"`
	ModifiedAt  time.Time         `json:"modifiedAt,omitempty"`
}

// SessionSpec holds session specification
type SessionSpec struct {
	DeploymentTargetID            string                        `json:"deploymentTargetId,omitempty"`
	FlinkConfiguration            map[string]string             `json:"flinkConfiguration,omitempty"`
	FlinkVersion                  string                        `json:"flinkVersion,omitempty"`
	SessionClusterResourceProfile SessionClusterResourceProfile `json:"sessionClusterResourceProfile,omitempty"`
}

// SessionClusterResourceProfile defines resource profile for session cluster
type SessionClusterResourceProfile struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// SessionStatus holds session status
type SessionStatus struct {
	State string `json:"state,omitempty"`
}

// SessionList represents a list of sessions
type SessionList struct {
	Items []Session `json:"items"`
}

// ListSessions lists all sessions in a namespace
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
func (c *Client) DeleteSession(namespace, name string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/v1/namespaces/%s/sessions/%s", namespace, name))

	return handleResponse(resp, err)
}
