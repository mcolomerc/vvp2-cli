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

var jobCmd = &cobra.Command{
	Use:     "job",
	Aliases: []string{"jobs"},
	Short:   "Manage VVP jobs",
	Long:    `List and view Ververica Platform jobs. Jobs represent running Flink applications.`,
}

var jobListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List jobs in a namespace",
	RunE:    runJobList,
}

var jobGetCmd = &cobra.Command{
	Use:   "get [jobId]",
	Short: "Get a job by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runJobGet,
}

func init() {
	rootCmd.AddCommand(jobCmd)

	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobGetCmd)

	// Add flags
	jobListCmd.Flags().StringP("namespace", "n", "", "Namespace")
	jobGetCmd.Flags().StringP("namespace", "n", "", "Namespace")
}

func runJobList(cmd *cobra.Command, args []string) error {
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

	jobs, err := client.ListJobs(namespace)
	if err != nil {
		return err
	}

	return printJobs(jobs.Items)
}

func runJobGet(cmd *cobra.Command, args []string) error {
	jobID := args[0]
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

	job, err := client.GetJob(namespace, jobID)
	if err != nil {
		return err
	}

	return printJob(job)
}

// Helper functions for printing jobs
func printJobs(jobs []api.Job) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(jobs, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(jobs)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format
		if len(jobs) == 0 {
			fmt.Println("No jobs found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "JOB ID\tNAME\tNAMESPACE\tSTATE\tDEPLOYMENT ID\tSTART TIME")
		for _, job := range jobs {
			jobID := job.Metadata.ID
			name := job.Metadata.Name
			if name == "" {
				name = "-"
			}
			state := job.Status.State
			deploymentID := job.Spec.DeploymentID
			
			startTime := "-"
			if job.Status.Running != nil && !job.Status.Running.StartTime.IsZero() {
				startTime = job.Status.Running.StartTime.Format("2006-01-02 15:04:05")
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				jobID,
				name,
				job.Metadata.Namespace,
				state,
				deploymentID,
				startTime,
			)
		}
		w.Flush()
	}
	return nil
}

func printJob(job *api.Job) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(job, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(job)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format with details
		fmt.Printf("Job ID: %s\n", job.Metadata.ID)
		if job.Metadata.Name != "" {
			fmt.Printf("Name: %s\n", job.Metadata.Name)
		}
		fmt.Printf("Namespace: %s\n", job.Metadata.Namespace)
		fmt.Printf("State: %s\n", job.Status.State)
		fmt.Printf("Deployment ID: %s\n", job.Spec.DeploymentID)

		if job.Status.Running != nil {
			fmt.Println("\nRunning Status:")
			if !job.Status.Running.StartTime.IsZero() {
				fmt.Printf("  Start Time: %s\n", job.Status.Running.StartTime.Format("2006-01-02 15:04:05"))
			}
			if !job.Status.Running.TransitionTime.IsZero() {
				fmt.Printf("  Transition Time: %s\n", job.Status.Running.TransitionTime.Format("2006-01-02 15:04:05"))
			}
			if job.Status.Running.JobID != "" {
				fmt.Printf("  Flink Job ID: %s\n", job.Status.Running.JobID)
			}
		}

		if job.Status.Failed != nil {
			fmt.Println("\nFailure Details:")
			if !job.Status.Failed.FailureTime.IsZero() {
				fmt.Printf("  Failure Time: %s\n", job.Status.Failed.FailureTime.Format("2006-01-02 15:04:05"))
			}
			if job.Status.Failed.Reason != "" {
				fmt.Printf("  Reason: %s\n", job.Status.Failed.Reason)
			}
			if job.Status.Failed.Message != "" {
				fmt.Printf("  Message: %s\n", job.Status.Failed.Message)
			}
		}

		if job.Status.Finished != nil && !job.Status.Finished.CompletionTime.IsZero() {
			fmt.Printf("\nCompletion Time: %s\n", job.Status.Finished.CompletionTime.Format("2006-01-02 15:04:05"))
		}

		if job.Status.Cancelled != nil && !job.Status.Cancelled.CancellationTime.IsZero() {
			fmt.Printf("\nCancellation Time: %s\n", job.Status.Cancelled.CancellationTime.Format("2006-01-02 15:04:05"))
		}

		if job.Status.Suspended != nil && !job.Status.Suspended.SuspensionTime.IsZero() {
			fmt.Printf("\nSuspension Time: %s\n", job.Status.Suspended.SuspensionTime.Format("2006-01-02 15:04:05"))
		}

		if len(job.Metadata.Labels) > 0 {
			fmt.Println("\nLabels:")
			for k, v := range job.Metadata.Labels {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}

		if len(job.Metadata.Annotations) > 0 {
			fmt.Println("\nAnnotations:")
			for k, v := range job.Metadata.Annotations {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}

		if !job.Metadata.CreatedAt.IsZero() {
			fmt.Printf("\nCreated At: %s\n", job.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		if !job.Metadata.ModifiedAt.IsZero() {
			fmt.Printf("Modified At: %s\n", job.Metadata.ModifiedAt.Format("2006-01-02 15:04:05"))
		}
	}
	return nil
}
