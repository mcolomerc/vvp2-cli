package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"mcolomerc/vvp2cli/pkg/api"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	deploymentDefaultsNamespace string
	deploymentDefaultsFile      string
)

// deploymentDefaultsCmd represents the deployment-defaults command group
var deploymentDefaultsCmd = &cobra.Command{
	Use:     "deployment-defaults",
	Aliases: []string{"defaults", "dd"},
	Short:   "Manage namespace deployment defaults",
	Long:    `Get and modify the namespace-level deployment defaults in Ververica Platform (Application Manager).`,
}

// getDeploymentDefaultsCmd fetches the deployment defaults
var getDeploymentDefaultsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get deployment defaults for a namespace",
	RunE:  runGetDeploymentDefaults,
}

// replaceDeploymentDefaultsCmd replaces the deployment defaults
var replaceDeploymentDefaultsCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace deployment defaults from a YAML/JSON file",
	RunE:  runReplaceDeploymentDefaults,
}

// updateDeploymentDefaultsCmd updates the deployment defaults via PATCH
var updateDeploymentDefaultsCmd = &cobra.Command{
	Use:   "update",
	Short: "Update deployment defaults via PATCH with a SecretValue YAML/JSON file",
	RunE:  runUpdateDeploymentDefaults,
}

func init() {
	rootCmd.AddCommand(deploymentDefaultsCmd)
	deploymentDefaultsCmd.AddCommand(getDeploymentDefaultsCmd)
	deploymentDefaultsCmd.AddCommand(replaceDeploymentDefaultsCmd)
	deploymentDefaultsCmd.AddCommand(updateDeploymentDefaultsCmd)

	// Flags
	deploymentDefaultsCmd.PersistentFlags().StringVarP(&deploymentDefaultsNamespace, "namespace", "n", "", "Namespace (defaults to config if not set)")

	replaceDeploymentDefaultsCmd.Flags().StringVarP(&deploymentDefaultsFile, "file", "f", "", "Path to deployment defaults YAML/JSON file (required)")
	replaceDeploymentDefaultsCmd.MarkFlagRequired("file")

	updateDeploymentDefaultsCmd.Flags().StringVarP(&deploymentDefaultsFile, "file", "f", "", "Path to SecretValue YAML/JSON file (required)")
	updateDeploymentDefaultsCmd.MarkFlagRequired("file")
}

func runGetDeploymentDefaults(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentDefaultsNamespace()
	if err != nil {
		return err
	}

	dd, err := client.GetDeploymentDefaults(ns)
	if err != nil {
		return fmt.Errorf("failed to get deployment defaults: %w", err)
	}

	return printDeploymentDefaults(dd)
}

func runReplaceDeploymentDefaults(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentDefaultsNamespace()
	if err != nil {
		return err
	}

	dd, err := loadDeploymentDefaultsFromFile(deploymentDefaultsFile)
	if err != nil {
		return err
	}

	res, err := client.ReplaceDeploymentDefaults(ns, dd)
	if err != nil {
		return fmt.Errorf("failed to replace deployment defaults: %w", err)
	}

	fmt.Println("Deployment defaults replaced successfully")
	return printDeploymentDefaults(res)
}

func runUpdateDeploymentDefaults(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentDefaultsNamespace()
	if err != nil {
		return err
	}

	sv, err := loadSecretValueFromFile(deploymentDefaultsFile)
	if err != nil {
		return err
	}

	res, err := client.UpdateDeploymentDefaults(ns, sv)
	if err != nil {
		return fmt.Errorf("failed to update deployment defaults: %w", err)
	}

	fmt.Println("Deployment defaults updated successfully")
	return printDeploymentDefaults(res)
}

func loadDeploymentDefaultsFromFile(filename string) (*api.DeploymentDefaults, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var dd api.DeploymentDefaults
	if err := json.Unmarshal(data, &dd); err != nil {
		if err := yaml.Unmarshal(data, &dd); err != nil {
			return nil, fmt.Errorf("failed to parse file as JSON or YAML: %w", err)
		}
	}
	return &dd, nil
}

func loadSecretValueFromFile(filename string) (*api.SecretValue, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var sv api.SecretValue
	if err := json.Unmarshal(data, &sv); err != nil {
		if err := yaml.Unmarshal(data, &sv); err != nil {
			return nil, fmt.Errorf("failed to parse file as JSON or YAML: %w", err)
		}
	}
	return &sv, nil
}

func printDeploymentDefaults(dd *api.DeploymentDefaults) error {
	format := GetConfig().GetOutputFormat()
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(dd)
	case "yaml":
		enc := yaml.NewEncoder(os.Stdout)
		enc.SetIndent(2)
		return enc.Encode(dd)
	default:
		// For table, print YAML for rich structure
		enc := yaml.NewEncoder(os.Stdout)
		enc.SetIndent(2)
		return enc.Encode(dd)
	}
}

// effectiveDeploymentDefaultsNamespace determines the namespace to use for deployment-defaults commands
func effectiveDeploymentDefaultsNamespace() (string, error) {
	if deploymentDefaultsNamespace != "" {
		return deploymentDefaultsNamespace, nil
	}
	if ns := GetConfig().GetNamespace(); ns != "" {
		return ns, nil
	}
	return "", fmt.Errorf("namespace not specified. Provide --namespace or set default.namespace in ~/.vvp2/config.yaml (or VVP_DEFAULT_NAMESPACE)")
}
