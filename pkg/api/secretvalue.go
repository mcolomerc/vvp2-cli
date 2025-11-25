package api

import (
	"fmt"
)

// SecretValueList represents a list of secret values
type SecretValueList struct {
	APIVersion string        `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string        `json:"kind,omitempty" yaml:"kind,omitempty"`
	Items      []SecretValue `json:"items" yaml:"items"`
}

// ListSecretValues lists all secret values in a namespace
func (c *Client) ListSecretValues(namespace string) (*SecretValueList, error) {
	var result SecretValueList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/secret-values", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSecretValue gets a secret value by name
func (c *Client) GetSecretValue(namespace, name string) (*SecretValue, error) {
	var result SecretValue
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/secret-values/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateSecretValue creates a new secret value
func (c *Client) CreateSecretValue(namespace string, secretValue *SecretValue) (*SecretValue, error) {
	var result SecretValue
	resp, err := c.httpClient.R().
		SetBody(secretValue).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/namespaces/%s/secret-values", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateSecretValue updates an existing secret value (PUT)
func (c *Client) UpdateSecretValue(namespace, name string, secretValue *SecretValue) (*SecretValue, error) {
	var result SecretValue
	resp, err := c.httpClient.R().
		SetBody(secretValue).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/namespaces/%s/secret-values/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteSecretValue deletes a secret value
func (c *Client) DeleteSecretValue(namespace, name string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/v1/namespaces/%s/secret-values/%s", namespace, name))

	return handleResponse(resp, err)
}

