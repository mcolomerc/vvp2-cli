package api

import (
	"crypto/tls"
	"fmt"
	"time"

	"mcolomerc/vvp2cli/pkg/config"

	"github.com/go-resty/resty/v2"
)

// Client is the VVP API client
type Client struct {
	httpClient *resty.Client
	baseURL    string
	token      string
}

// NewClient creates a new VVP API client
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	httpClient := resty.New()
	httpClient.SetTimeout(30 * time.Second)
	httpClient.SetBaseURL(cfg.GetAPIURL())

	// Set token if provided
	if cfg.GetToken() != "" {
		httpClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", cfg.GetToken()))
	}

	// Skip TLS verification if configured
	if cfg.IsInsecure() {
		httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    cfg.GetAPIURL(),
		token:      cfg.GetToken(),
	}, nil
}

// SetDebug enables debug mode for the HTTP client
func (c *Client) SetDebug(debug bool) {
	c.httpClient.SetDebug(debug)
}

// APIError represents an API error response
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// handleResponse checks the response and returns an error if needed
func handleResponse(resp *resty.Response, err error) error {
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return &APIError{
			StatusCode: resp.StatusCode(),
			Message:    string(resp.Body()),
		}
	}

	return nil
}
