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
	deploymentNamespace string
	deploymentFile      string
	deploymentState     string
)

// deploymentCmd represents the deployment command
var deploymentCmd = &cobra.Command{
	Use:     "deployment",
	Aliases: []string{"deploy", "deployments"},
	Short:   "Manage VVP deployments",
	Long:    `Create, read, update, and delete Ververica Platform deployments.`,
}

// listDeploymentsCmd lists all deployments
var listDeploymentsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all deployments",
	RunE:    runListDeployments,
}

// getDeploymentCmd gets a deployment by name
var getDeploymentCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a deployment by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetDeployment,
}

// createDeploymentCmd creates a new deployment
var createDeploymentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new deployment from a YAML/JSON file",
	RunE:  runCreateDeployment,
}

// updateDeploymentCmd updates a deployment
var updateDeploymentCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update an existing deployment",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdateDeployment,
}

// deleteDeploymentCmd deletes a deployment
var deleteDeploymentCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a deployment",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeleteDeployment,
}

// startDeploymentCmd starts a deployment
var startDeploymentCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Start a deployment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdateDeploymentState(args[0], "RUNNING")
	},
}

// stopDeploymentCmd stops a deployment
var stopDeploymentCmd = &cobra.Command{
	Use:   "stop [name]",
	Short: "Stop a deployment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdateDeploymentState(args[0], "CANCELLED")
	},
}

// suspendDeploymentCmd suspends a deployment
var suspendDeploymentCmd = &cobra.Command{
	Use:   "suspend [name]",
	Short: "Suspend a deployment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdateDeploymentState(args[0], "SUSPENDED")
	},
}

func init() {
	rootCmd.AddCommand(deploymentCmd)
	deploymentCmd.AddCommand(listDeploymentsCmd)
	deploymentCmd.AddCommand(getDeploymentCmd)
	deploymentCmd.AddCommand(createDeploymentCmd)
	deploymentCmd.AddCommand(updateDeploymentCmd)
	deploymentCmd.AddCommand(deleteDeploymentCmd)
	deploymentCmd.AddCommand(startDeploymentCmd)
	deploymentCmd.AddCommand(stopDeploymentCmd)
	deploymentCmd.AddCommand(suspendDeploymentCmd)

	// Flags for deployment commands
	deploymentCmd.PersistentFlags().StringVarP(&deploymentNamespace, "namespace", "n", "", "Namespace (defaults to config if not set)")

	createDeploymentCmd.Flags().StringVarP(&deploymentFile, "file", "f", "", "Path to deployment YAML/JSON file (required)")
	createDeploymentCmd.MarkFlagRequired("file")

	updateDeploymentCmd.Flags().StringVarP(&deploymentFile, "file", "f", "", "Path to deployment YAML/JSON file (required)")
	updateDeploymentCmd.MarkFlagRequired("file")
	
	deleteDeploymentCmd.Flags().BoolP("force", "", false, "Force delete by cancelling the deployment first if needed")
	updateDeploymentCmd.MarkFlagRequired("file")
}

func runListDeployments(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentNamespace()
	if err != nil {
		return err
	}

	deployments, err := client.ListDeployments(ns)
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	return printDeployments(deployments.Items)
}

func runGetDeployment(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentNamespace()
	if err != nil {
		return err
	}

	deployment, err := client.GetDeployment(ns, args[0])
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	return printDeployment(deployment)
}

func runCreateDeployment(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentNamespace()
	if err != nil {
		return err
	}

	deployment, err := loadDeploymentFromFile(deploymentFile)
	if err != nil {
		return err
	}

	result, err := client.CreateDeployment(ns, deployment)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	fmt.Printf("Deployment %s created successfully\n", result.Metadata.Name)
	return printDeployment(result)
}

func runUpdateDeployment(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentNamespace()
	if err != nil {
		return err
	}

	deployment, err := loadDeploymentFromFile(deploymentFile)
	if err != nil {
		return err
	}

	result, err := client.UpdateDeployment(ns, args[0], deployment)
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	fmt.Printf("Deployment %s updated successfully\n", result.Metadata.Name)
	return printDeployment(result)
}

func runDeleteDeployment(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentNamespace()
	if err != nil {
		return err
	}

	force, _ := cmd.Flags().GetBool("force")
	
	// If force flag is set, try to cancel the deployment first
	if force {
		// Get current deployment state
		deployment, err := client.GetDeployment(ns, args[0])
		if err != nil {
			return fmt.Errorf("failed to get deployment: %w", err)
		}
		
		// If not already cancelled, cancel it first
		if deployment.Spec.State != "CANCELLED" {
			fmt.Printf("Cancelling deployment %s before deletion...\n", args[0])
			if _, err := client.UpdateDeploymentState(ns, args[0], "CANCELLED"); err != nil {
				return fmt.Errorf("failed to cancel deployment: %w", err)
			}
			fmt.Printf("Deployment %s transitioned to CANCELLED state\n", args[0])
			fmt.Printf("Note: The deployment may take some time to fully cancel. If deletion fails, wait a moment and try again.\n")
		}
	}

	if err := client.DeleteDeployment(ns, args[0]); err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	fmt.Printf("Deployment %s deleted successfully\n", args[0])
	return nil
}

func runUpdateDeploymentState(name, state string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveDeploymentNamespace()
	if err != nil {
		return err
	}

	result, err := client.UpdateDeploymentState(ns, name, state)
	if err != nil {
		return fmt.Errorf("failed to update deployment state: %w", err)
	}

	fmt.Printf("Deployment %s state updated to %s\n", name, state)
	return printDeployment(result)
}

func loadDeploymentFromFile(filename string) (*api.Deployment, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var deployment api.Deployment

	// Try JSON first
	if err := json.Unmarshal(data, &deployment); err != nil {
		// Try YAML
		if err := yaml.Unmarshal(data, &deployment); err != nil {
			return nil, fmt.Errorf("failed to parse file as JSON or YAML: %w", err)
		}
	}

	return &deployment, nil
}

func printDeployments(deployments []api.Deployment) error {
	if len(deployments) == 0 {
		fmt.Println("No deployments found")
		return nil
	}

	format := GetConfig().GetOutputFormat()

	switch format {
	case "json":
		return printJSON(deployments)
	case "yaml":
		return printYAML(deployments)
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tNAMESPACE\tSTATE\tCREATED")
		for _, d := range deployments {
			state := "N/A"
			if d.Status != nil && d.Status.State != "" {
				state = d.Status.State
			} else if d.Spec.State != "" {
				state = d.Spec.State
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				d.Metadata.Name,
				d.Metadata.Namespace,
				state,
				d.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		return w.Flush()
	}
}

func printDeployment(deployment *api.Deployment) error {
	format := GetConfig().GetOutputFormat()

	switch format {
	case "json":
		return printJSON(deployment)
	case "yaml":
		return printYAML(deployment)
	default:
		return printYAML(deployment)
	}
}

func printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func printYAML(v interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(v)
}

// effectiveDeploymentNamespace returns the namespace for deployment commands,
// preferring the --namespace flag, then falling back to the default from config.
// If neither is set, it returns a helpful error.
func effectiveDeploymentNamespace() (string, error) {
	if deploymentNamespace != "" {
		return deploymentNamespace, nil
	}
	if ns := GetConfig().GetNamespace(); ns != "" {
		return ns, nil
	}
	return "", fmt.Errorf("namespace not specified. Provide --namespace or set default.namespace in ~/.vvp2/config.yaml (or VVP_DEFAULT_NAMESPACE)")
}
