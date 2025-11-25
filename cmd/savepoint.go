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

var savepointCmd = &cobra.Command{
	Use:     "savepoint",
	Aliases: []string{"savepoints", "sp"},
	Short:   "Manage VVP savepoints",
	Long:    `List, view, create, and delete Ververica Platform savepoints. Savepoints are snapshots of Flink job state.`,
}

var savepointListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List savepoints in a namespace",
	RunE:    runSavepointList,
}

var savepointGetCmd = &cobra.Command{
	Use:   "get [savepointId]",
	Short: "Get a savepoint by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runSavepointGet,
}

var savepointCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a savepoint for a deployment or job",
	Long: `Create a savepoint for a deployment or job.
You must specify either --deployment-id or --job-id.`,
	RunE: runSavepointCreate,
}

var savepointDeleteCmd = &cobra.Command{
	Use:     "delete [savepointId]",
	Aliases: []string{"rm"},
	Short:   "Delete a savepoint",
	Args:    cobra.ExactArgs(1),
	RunE:    runSavepointDelete,
}

func init() {
	rootCmd.AddCommand(savepointCmd)

	savepointCmd.AddCommand(savepointListCmd)
	savepointCmd.AddCommand(savepointGetCmd)
	savepointCmd.AddCommand(savepointCreateCmd)
	savepointCmd.AddCommand(savepointDeleteCmd)

	// Add flags
	savepointListCmd.Flags().StringP("namespace", "n", "", "Namespace")
	savepointGetCmd.Flags().StringP("namespace", "n", "", "Namespace")
	
	savepointCreateCmd.Flags().StringP("namespace", "n", "", "Namespace")
	savepointCreateCmd.Flags().String("deployment-id", "", "Deployment ID to create savepoint for")
	savepointCreateCmd.Flags().String("job-id", "", "Job ID to create savepoint for")
	savepointCreateCmd.Flags().String("name", "", "Optional name for the savepoint")
	
	savepointDeleteCmd.Flags().StringP("namespace", "n", "", "Namespace")
}

func runSavepointList(cmd *cobra.Command, args []string) error {
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = cfg.Default.Namespace
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	savepoints, err := client.ListSavepoints(namespace)
	if err != nil {
		return err
	}

	return printSavepoints(savepoints.Items)
}

func runSavepointGet(cmd *cobra.Command, args []string) error {
	savepointID := args[0]
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = cfg.Default.Namespace
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	savepoint, err := client.GetSavepoint(namespace, savepointID)
	if err != nil {
		return err
	}

	return printSavepoint(savepoint)
}

func runSavepointCreate(cmd *cobra.Command, args []string) error {
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = cfg.Default.Namespace
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	deploymentID, _ := cmd.Flags().GetString("deployment-id")
	jobID, _ := cmd.Flags().GetString("job-id")
	name, _ := cmd.Flags().GetString("name")

	if deploymentID == "" && jobID == "" {
		return fmt.Errorf("either --deployment-id or --job-id must be specified")
	}

	if deploymentID != "" && jobID != "" {
		return fmt.Errorf("cannot specify both --deployment-id and --job-id")
	}

	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	request := &api.SavepointCreationRequest{
		Metadata: api.SavepointMetadata{
			Namespace: namespace,
		},
		Spec: api.SavepointSpec{
			DeploymentID: deploymentID,
			JobID:        jobID,
		},
	}

	if name != "" {
		request.Metadata.Name = name
	}

	result, err := client.CreateSavepoint(namespace, request)
	if err != nil {
		return err
	}

	fmt.Printf("Savepoint creation initiated: %s\n", result.Metadata.ID)
	return printSavepoint(result)
}

func runSavepointDelete(cmd *cobra.Command, args []string) error {
	savepointID := args[0]
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = cfg.Default.Namespace
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	if err := client.DeleteSavepoint(namespace, savepointID); err != nil {
		return err
	}

	fmt.Printf("Savepoint '%s' deleted successfully\n", savepointID)
	return nil
}

// Helper functions for printing savepoints
func printSavepoints(savepoints []api.Savepoint) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(savepoints, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(savepoints)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format
		if len(savepoints) == 0 {
			fmt.Println("No savepoints found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "SAVEPOINT ID\tNAME\tNAMESPACE\tSTATE\tDEPLOYMENT ID\tJOB ID\tCREATED")
		for _, sp := range savepoints {
			savepointID := sp.Metadata.ID
			name := sp.Metadata.Name
			if name == "" {
				name = "-"
			}
			state := sp.Status.State
			deploymentID := sp.Spec.DeploymentID
			if deploymentID == "" {
				deploymentID = "-"
			}
			jobID := sp.Spec.JobID
			if jobID == "" {
				jobID = "-"
			}
			
			created := "-"
			if !sp.Metadata.CreatedAt.IsZero() {
				created = sp.Metadata.CreatedAt.Format("2006-01-02 15:04:05")
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				savepointID,
				name,
				sp.Metadata.Namespace,
				state,
				deploymentID,
				jobID,
				created,
			)
		}
		w.Flush()
	}
	return nil
}

func printSavepoint(sp *api.Savepoint) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(sp, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(sp)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format with details
		fmt.Printf("Savepoint ID: %s\n", sp.Metadata.ID)
		if sp.Metadata.Name != "" {
			fmt.Printf("Name: %s\n", sp.Metadata.Name)
		}
		fmt.Printf("Namespace: %s\n", sp.Metadata.Namespace)
		fmt.Printf("State: %s\n", sp.Status.State)
		
		if sp.Spec.DeploymentID != "" {
			fmt.Printf("Deployment ID: %s\n", sp.Spec.DeploymentID)
		}
		if sp.Spec.JobID != "" {
			fmt.Printf("Job ID: %s\n", sp.Spec.JobID)
		}

		if sp.Status.Completed != nil {
			fmt.Println("\nCompleted Status:")
			if sp.Status.Completed.Location != "" {
				fmt.Printf("  Location: %s\n", sp.Status.Completed.Location)
			}
			if !sp.Status.Completed.Time.IsZero() {
				fmt.Printf("  Completion Time: %s\n", sp.Status.Completed.Time.Format("2006-01-02 15:04:05"))
			}
		}

		if sp.Status.Failed != nil {
			fmt.Println("\nFailure Details:")
			if !sp.Status.Failed.Time.IsZero() {
				fmt.Printf("  Failure Time: %s\n", sp.Status.Failed.Time.Format("2006-01-02 15:04:05"))
			}
			if sp.Status.Failed.Reason != "" {
				fmt.Printf("  Reason: %s\n", sp.Status.Failed.Reason)
			}
			if sp.Status.Failed.Message != "" {
				fmt.Printf("  Message: %s\n", sp.Status.Failed.Message)
			}
		}

		if len(sp.Metadata.Labels) > 0 {
			fmt.Println("\nLabels:")
			for k, v := range sp.Metadata.Labels {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}

		if len(sp.Metadata.Annotations) > 0 {
			fmt.Println("\nAnnotations:")
			for k, v := range sp.Metadata.Annotations {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}

		if !sp.Metadata.CreatedAt.IsZero() {
			fmt.Printf("\nCreated At: %s\n", sp.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		if !sp.Metadata.ModifiedAt.IsZero() {
			fmt.Printf("Modified At: %s\n", sp.Metadata.ModifiedAt.Format("2006-01-02 15:04:05"))
		}
	}
	return nil
}
