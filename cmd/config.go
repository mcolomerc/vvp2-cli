package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ConfigData represents the configuration structure
type ConfigData struct {
	API struct {
		URL      string `yaml:"url"`
		Token    string `yaml:"token"`
		Insecure bool   `yaml:"insecure"`
	} `yaml:"api"`
	Default struct {
		Namespace string `yaml:"namespace"`
	} `yaml:"default"`
	Output struct {
		Format string `yaml:"format"`
	} `yaml:"output"`
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage vvp2 configuration",
	Long:  `Manage vvp2 configuration file stored in ~/.vvp2/config.yaml`,
}

// configInitCmd initializes the configuration
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize vvp2 configuration interactively",
	Long:  `Creates a configuration file at ~/.vvp2/config.yaml by prompting for values.`,
	RunE:  runConfigInit,
}

// configShowCmd shows the current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration file contents.`,
	RunE:  runConfigShow,
}

// configPathCmd shows the configuration file path
var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  `Display the path to the configuration file.`,
	RunE:  runConfigPath,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)

	// Flags for config init
	configInitCmd.Flags().BoolP("force", "f", false, "Overwrite existing configuration file")
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create .vvp2 directory if it doesn't exist
	configDir := filepath.Join(home, ".vvp2")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil && !force {
		return fmt.Errorf("configuration file already exists at %s. Use --force to overwrite", configPath)
	}

	fmt.Println("VVP2 Configuration Setup")
	fmt.Println("========================")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	config := ConfigData{}

	// API URL
	fmt.Print("Ververica Platform API URL [http://vvp.localhost]: ")
	apiURL, _ := reader.ReadString('\n')
	apiURL = strings.TrimSpace(apiURL)
	if apiURL == "" {
		apiURL = "http://vvp.localhost"
	}
	config.API.URL = apiURL

	// API Token
	fmt.Print("API Token (leave empty if not required): ")
	apiToken, _ := reader.ReadString('\n')
	apiToken = strings.TrimSpace(apiToken)
	config.API.Token = apiToken

	// Insecure
	fmt.Print("Skip TLS certificate verification? [y/N]: ")
	insecureInput, _ := reader.ReadString('\n')
	insecureInput = strings.TrimSpace(strings.ToLower(insecureInput))
	config.API.Insecure = insecureInput == "y" || insecureInput == "yes"

	// Default namespace
	fmt.Print("Default namespace [default]: ")
	namespace, _ := reader.ReadString('\n')
	namespace = strings.TrimSpace(namespace)
	if namespace == "" {
		namespace = "default"
	}
	config.Default.Namespace = namespace

	// Output format
	fmt.Print("Default output format (table/json/yaml) [table]: ")
	format, _ := reader.ReadString('\n')
	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" {
		format = "table"
	}
	if format != "table" && format != "json" && format != "yaml" {
		format = "table"
	}
	config.Output.Format = format

	// Write configuration file
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	fmt.Println()
	fmt.Printf("✓ Configuration file created at: %s\n", configPath)
	fmt.Println()
	fmt.Println("You can now use vvp2 commands:")
	fmt.Println("  vvp2 namespace list")
	fmt.Println("  vvp2 deployment list -n", config.Default.Namespace)
	fmt.Println()
	fmt.Println("To view your configuration:")
	fmt.Println("  vvp2 config show")

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found at %s. Run 'vvp2 config init' to create it", configPath)
	}

	// Read and display configuration
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	fmt.Printf("Configuration file: %s\n", configPath)
	fmt.Println("---")
	fmt.Print(string(data))

	return nil
}

func runConfigPath(cmd *cobra.Command, args []string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Check if config file exists
	exists := "does not exist"
	if _, err := os.Stat(configPath); err == nil {
		exists = "exists"
	}

	fmt.Printf("Config file path: %s (%s)\n", configPath, exists)
	fmt.Printf("Config directory: %s\n", filepath.Dir(configPath))

	if exists == "does not exist" {
		fmt.Println()
		fmt.Println("Run 'vvp2 config init' to create the configuration file.")
	}

	return nil
}

func getConfigPath() (string, error) {
	if cfgFile != "" {
		return cfgFile, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".vvp2", "config.yaml"), nil
}
