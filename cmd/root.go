package cmd

import (
	"fmt"
	"os"

	"mcolomerc/vvp2cli/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	cfg       *config.Config
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

// SetVersionInfo sets the version information from the build
func SetVersionInfo(v, c, bt string) {
	version = v
	commit = c
	buildTime = bt
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vvp2",
	Short: "A CLI tool to interact with Ververica Platform API",
	Long: `vvp2 is a command-line interface tool for interacting with the Ververica Platform API.
It provides commands to manage deployments, namespaces, session clusters, and other VVP resources.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vvp2/config.yaml)")
	rootCmd.PersistentFlags().String("api-url", "", "Ververica Platform API URL")
	rootCmd.PersistentFlags().String("api-token", "", "API authentication token")
	rootCmd.PersistentFlags().String("namespace", "", "Default namespace")
	rootCmd.PersistentFlags().Bool("insecure", false, "Skip TLS certificate verification")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table, json, yaml)")

	// Bind flags to viper
	viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("api.token", rootCmd.PersistentFlags().Lookup("api-token"))
	viper.BindPFlag("api.insecure", rootCmd.PersistentFlags().Lookup("insecure"))
	viper.BindPFlag("default.namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("output.format", rootCmd.PersistentFlags().Lookup("output"))

	// Add usage command
	rootCmd.AddCommand(usageCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory under .vvp2 folder
		configDir := home + "/.vvp2"
		viper.AddConfigPath(configDir)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Environment variables
	viper.SetEnvPrefix("VVP")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Load configuration
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load configuration: %v\n", err)
	}
}

// GetConfig returns the current configuration
func GetConfig() *config.Config {
	return cfg
}
