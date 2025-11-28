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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get VVP platform status",
	Long:  `Retrieve the status of the Ververica Platform, including health, version, and component information.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	status, err := client.GetStatus()
	if err != nil {
		return err
	}

	return printStatus(status)
}

func printStatus(status *api.Status) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(status)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format with details
		fmt.Println("=== Platform Status ===")

		// Health Status
		fmt.Println("\nHealth:")
		if status.Health.Status != "" {
			fmt.Printf("  Status: %s\n", status.Health.Status)
		}
		if status.Health.Message != "" {
			fmt.Printf("  Message: %s\n", status.Health.Message)
		}

		// Version Information
		if status.Version.Platform != "" || status.Version.Flink != "" {
			fmt.Println("\nVersion:")
			if status.Version.Platform != "" {
				fmt.Printf("  Platform: %s\n", status.Version.Platform)
			}
			if status.Version.Edition != "" {
				fmt.Printf("  Edition: %s\n", status.Version.Edition)
			}
			if status.Version.Flink != "" {
				fmt.Printf("  Flink: %s\n", status.Version.Flink)
			}
			if status.Version.BuildTime != "" {
				fmt.Printf("  Build Time: %s\n", status.Version.BuildTime)
			}
			if status.Version.CommitHash != "" {
				fmt.Printf("  Commit: %s\n", status.Version.CommitHash)
			}
		}

		// Components
		if len(status.Components) > 0 {
			fmt.Println("\nComponents:")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "  NAME\tSTATUS\tVERSION\tMESSAGE")
			for _, component := range status.Components {
				name := component.Name
				compStatus := component.Status
				version := component.Version
				if version == "" {
					version = "-"
				}
				message := component.Message
				if message == "" {
					message = "-"
				}
				fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", name, compStatus, version, message)
			}
			w.Flush()
		}

		// Resource Usage
		if status.ResourceUsage != nil {
			fmt.Println("\nResource Usage:")
			if status.ResourceUsage.Namespaces > 0 {
				fmt.Printf("  Namespaces: %d\n", status.ResourceUsage.Namespaces)
			}
			if status.ResourceUsage.Deployments > 0 {
				fmt.Printf("  Deployments: %d\n", status.ResourceUsage.Deployments)
			}
			if status.ResourceUsage.Jobs > 0 {
				fmt.Printf("  Jobs: %d\n", status.ResourceUsage.Jobs)
			}
			if status.ResourceUsage.SessionClusters > 0 {
				fmt.Printf("  Session Clusters: %d\n", status.ResourceUsage.SessionClusters)
			}
		}
	}
	return nil
}
