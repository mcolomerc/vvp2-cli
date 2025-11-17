package api

import (
	"fmt"
	"time"
)

// Namespace represents a VVP namespace
type Namespace struct {
	Metadata NamespaceMetadata `json:"metadata"`
	Spec     NamespaceSpec     `json:"spec,omitempty"`
	Status   NamespaceStatus   `json:"status,omitempty"`
}

// NamespaceMetadata holds namespace metadata
type NamespaceMetadata struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"createdAt,omitempty"`
	ModifiedAt  time.Time         `json:"modifiedAt,omitempty"`
}

// NamespaceSpec holds namespace specification
type NamespaceSpec struct {
	RoleBindings []RoleBinding `json:"roleBindings,omitempty"`
}

// NamespaceStatus holds namespace status
type NamespaceStatus struct {
	State string `json:"state,omitempty"`
}

// RoleBinding represents a role binding
type RoleBinding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}

// NamespaceList represents a list of namespaces
type NamespaceList struct {
	Items []Namespace `json:"items"`
}

// ListNamespaces lists all namespaces
func (c *Client) ListNamespaces() (*NamespaceList, error) {
	var result NamespaceList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get("/namespaces/v1/namespaces")

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetNamespace gets a namespace by name
func (c *Client) GetNamespace(name string) (*Namespace, error) {
	var result Namespace
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/namespaces/v1/namespaces/%s", name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateNamespace creates a new namespace
func (c *Client) CreateNamespace(namespace *Namespace) (*Namespace, error) {
	var result Namespace
	resp, err := c.httpClient.R().
		SetBody(namespace).
		SetResult(&result).
		Post("/namespaces/v1/namespaces")

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateNamespace updates an existing namespace
func (c *Client) UpdateNamespace(name string, namespace *Namespace) (*Namespace, error) {
	var result Namespace
	resp, err := c.httpClient.R().
		SetBody(namespace).
		SetResult(&result).
		Put(fmt.Sprintf("/namespaces/v1/namespaces/%s", name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteNamespace deletes a namespace
func (c *Client) DeleteNamespace(name string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/namespaces/v1/namespaces/%s", name))

	return handleResponse(resp, err)
}
