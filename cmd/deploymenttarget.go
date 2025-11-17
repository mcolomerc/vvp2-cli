package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"mcolomerc/vvp2cli/pkg/api"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	deploymentTargetNamespace string
	deploymentTargetFile      string
)

// deploymentTargetCmd represents the deployment-target command
var deploymentTargetCmd = &cobra.Command{
	Use:     "deployment-target",
	Aliases: []string{"dt", "deployment-targets", "target", "targets"},
	Short:   "Manage VVP deployment targets",
	Long:    `Create, read, update, and delete Ververica Platform deployment targets.`,
}

// listDeploymentTargetsCmd lists all deployment targets
var listDeploymentTargetsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all deployment targets",
	RunE:    runListDeploymentTargets,
}

// getDeploymentTargetCmd gets a deployment target by name
var getDeploymentTargetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a deployment target by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetDeploymentTarget,
}

// createDeploymentTargetCmd creates a new deployment target
var createDeploymentTargetCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new deployment target from a YAML/JSON file",
	RunE:  runCreateDeploymentTarget,
}

// updateDeploymentTargetCmd updates a deployment target
var updateDeploymentTargetCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update an existing deployment target",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdateDeploymentTarget,
}

// deleteDeploymentTargetCmd deletes a deployment target
var deleteDeploymentTargetCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a deployment target",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeleteDeploymentTarget,
}

func init() {
	rootCmd.AddCommand(deploymentTargetCmd)
	deploymentTargetCmd.AddCommand(listDeploymentTargetsCmd)
	deploymentTargetCmd.AddCommand(getDeploymentTargetCmd)
	deploymentTargetCmd.AddCommand(createDeploymentTargetCmd)
	deploymentTargetCmd.AddCommand(updateDeploymentTargetCmd)
	deploymentTargetCmd.AddCommand(deleteDeploymentTargetCmd)

	// Flags for deployment target commands
	deploymentTargetCmd.PersistentFlags().StringVarP(&deploymentTargetNamespace, "namespace", "n", "", "Namespace (defaults to config if not set)")

	createDeploymentTargetCmd.Flags().StringVarP(&deploymentTargetFile, "file", "f", "", "Path to deployment target YAML/JSON file (required)")
	createDeploymentTargetCmd.MarkFlagRequired("file")

	updateDeploymentTargetCmd.Flags().StringVarP(&deploymentTargetFile, "file", "f", "", "Path to deployment target YAML/JSON file (required)")
	updateDeploymentTargetCmd.MarkFlagRequired("file")
}

func runListDeploymentTargets(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentTargetNamespace()
	if err != nil {
		return err
	}

	targets, err := client.ListDeploymentTargets(ns)
	if err != nil {
		return fmt.Errorf("failed to list deployment targets: %w", err)
	}

	return printDeploymentTargets(targets.Items)
}

func runGetDeploymentTarget(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentTargetNamespace()
	if err != nil {
		return err
	}

	target, err := client.GetDeploymentTarget(ns, args[0])
	if err != nil {
		return fmt.Errorf("failed to get deployment target: %w", err)
	}

	return printDeploymentTarget(target)
}

func runCreateDeploymentTarget(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentTargetNamespace()
	if err != nil {
		return err
	}

	target, err := loadDeploymentTargetFromFile(deploymentTargetFile)
	if err != nil {
		return err
	}

	result, err := client.CreateDeploymentTarget(ns, target)
	if err != nil {
		return fmt.Errorf("failed to create deployment target: %w", err)
	}

	fmt.Printf("Deployment target %s created successfully\n", result.Metadata.Name)
	return printDeploymentTarget(result)
}

func runUpdateDeploymentTarget(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentTargetNamespace()
	if err != nil {
		return err
	}

	target, err := loadDeploymentTargetFromFile(deploymentTargetFile)
	if err != nil {
		return err
	}

	result, err := client.UpdateDeploymentTarget(ns, args[0], target)
	if err != nil {
		return fmt.Errorf("failed to update deployment target: %w", err)
	}

	fmt.Printf("Deployment target %s updated successfully\n", result.Metadata.Name)
	return printDeploymentTarget(result)
}

func runDeleteDeploymentTarget(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentTargetNamespace()
	if err != nil {
		return err
	}

	if err := client.DeleteDeploymentTarget(ns, args[0]); err != nil {
		return fmt.Errorf("failed to delete deployment target: %w", err)
	}

	fmt.Printf("Deployment target %s deleted successfully\n", args[0])
	return nil
}

func loadDeploymentTargetFromFile(filename string) (*api.DeploymentTargetResource, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var target api.DeploymentTargetResource

	// Try JSON first
	if err := json.Unmarshal(data, &target); err != nil {
		// Try YAML
		if err := yaml.Unmarshal(data, &target); err != nil {
			return nil, fmt.Errorf("failed to parse file as JSON or YAML: %w", err)
		}
	}

	return &target, nil
}

func printDeploymentTargets(targets []api.DeploymentTargetResource) error {
	if len(targets) == 0 {
		fmt.Println("No deployment targets found")
		return nil
	}

	format := GetConfig().GetOutputFormat()

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(targets)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(targets)
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tNAMESPACE\tSTATE\tCREATED")
		for _, t := range targets {
			state := t.Status.State
			if state == "" {
				state = "N/A"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				t.Metadata.Name,
				t.Metadata.Namespace,
				state,
				t.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		return w.Flush()
	}
}

func printDeploymentTarget(target *api.DeploymentTargetResource) error {
	format := GetConfig().GetOutputFormat()

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(target)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(target)
	default:
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(target)
	}
}

// effectiveDeploymentTargetNamespace returns the namespace for deployment-target commands,
// preferring the --namespace flag, then falling back to the default from config.
// If neither is set, it returns a helpful error.
func effectiveDeploymentTargetNamespace() (string, error) {
	if deploymentTargetNamespace != "" {
		return deploymentTargetNamespace, nil
	}
	if ns := GetConfig().GetNamespace(); ns != "" {
		return ns, nil
	}
	return "", fmt.Errorf("namespace not specified. Provide --namespace or set default.namespace in ~/.vvp2/config.yaml (or VVP_DEFAULT_NAMESPACE)")
}
