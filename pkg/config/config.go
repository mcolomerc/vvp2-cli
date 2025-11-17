package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	API     APIConfig     `mapstructure:"api"`
	Default DefaultConfig `mapstructure:"default"`
	Output  OutputConfig  `mapstructure:"output"`
}

// APIConfig holds API-related configuration
type APIConfig struct {
	URL      string `mapstructure:"url"`
	Token    string `mapstructure:"token"`
	Insecure bool   `mapstructure:"insecure"`
}

// DefaultConfig holds default values
type DefaultConfig struct {
	Namespace string `mapstructure:"namespace"`
}

// OutputConfig holds output formatting configuration
type OutputConfig struct {
	Format string `mapstructure:"format"`
}

// LoadConfig loads configuration from viper
func LoadConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.API.URL == "" {
		return fmt.Errorf("API URL is required (set via --api-url flag, VVP_API_URL env var, or config file)")
	}
	return nil
}

// GetAPIURL returns the configured API URL
func (c *Config) GetAPIURL() string {
	return c.API.URL
}

// GetToken returns the configured API token
func (c *Config) GetToken() string {
	return c.API.Token
}

// GetNamespace returns the default namespace
func (c *Config) GetNamespace() string {
	return c.Default.Namespace
}

// IsInsecure returns whether to skip TLS verification
func (c *Config) IsInsecure() bool {
	return c.API.Insecure
}

// GetOutputFormat returns the output format
func (c *Config) GetOutputFormat() string {
	if c.Output.Format == "" {
		return "table"
	}
	return c.Output.Format
}
