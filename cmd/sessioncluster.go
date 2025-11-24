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

var sessionClusterCmd = &cobra.Command{
	Use:     "sessioncluster",
	Aliases: []string{"sc", "session-cluster"},
	Short:   "Manage VVP session clusters",
	Long:    `Manage Ververica Platform session clusters (SQL session clusters).`,
}

var sessionClusterListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List session clusters in a namespace",
	RunE:    runSessionClusterList,
}

var sessionClusterGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a session cluster by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runSessionClusterGet,
}

var sessionClusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a session cluster from a file",
	RunE:  runSessionClusterCreate,
}

var sessionClusterUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update a session cluster",
	Args:  cobra.ExactArgs(1),
	RunE:  runSessionClusterUpdate,
}

var sessionClusterDeleteCmd = &cobra.Command{
	Use:     "delete [name]",
	Aliases: []string{"rm"},
	Short:   "Delete a session cluster",
	Args:    cobra.ExactArgs(1),
	RunE:    runSessionClusterDelete,
}

func init() {
	rootCmd.AddCommand(sessionClusterCmd)

	sessionClusterCmd.AddCommand(sessionClusterListCmd)
	sessionClusterCmd.AddCommand(sessionClusterGetCmd)
	sessionClusterCmd.AddCommand(sessionClusterCreateCmd)
	sessionClusterCmd.AddCommand(sessionClusterUpdateCmd)
	sessionClusterCmd.AddCommand(sessionClusterDeleteCmd)

	// Add flags
	sessionClusterListCmd.Flags().StringP("namespace", "n", "", "Namespace")
	sessionClusterGetCmd.Flags().StringP("namespace", "n", "", "Namespace")
	sessionClusterCreateCmd.Flags().StringP("namespace", "n", "", "Namespace")
	sessionClusterCreateCmd.Flags().StringP("file", "f", "", "File containing session cluster definition")
	sessionClusterCreateCmd.MarkFlagRequired("file")
	sessionClusterUpdateCmd.Flags().StringP("namespace", "n", "", "Namespace")
	sessionClusterUpdateCmd.Flags().StringP("file", "f", "", "File containing session cluster definition")
	sessionClusterUpdateCmd.MarkFlagRequired("file")
	sessionClusterDeleteCmd.Flags().StringP("namespace", "n", "", "Namespace")
}

func runSessionClusterList(cmd *cobra.Command, args []string) error {
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

	sessionClusters, err := client.ListSessionClusters(namespace)
	if err != nil {
		return err
	}

	return printSessionClusters(sessionClusters.Items)
}

func runSessionClusterGet(cmd *cobra.Command, args []string) error {
	name := args[0]
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

	sessionCluster, err := client.GetSessionCluster(namespace, name)
	if err != nil {
		return err
	}

	return printSessionCluster(sessionCluster)
}

func runSessionClusterCreate(cmd *cobra.Command, args []string) error {
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = cfg.Default.Namespace
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	filename, _ := cmd.Flags().GetString("file")
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var sessionCluster api.SessionCluster
	if err := yaml.Unmarshal(data, &sessionCluster); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Set namespace from flag if not in file
	if sessionCluster.Metadata.Namespace == "" {
		sessionCluster.Metadata.Namespace = namespace
	}

	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	result, err := client.CreateSessionCluster(namespace, &sessionCluster)
	if err != nil {
		return err
	}

	fmt.Printf("Session cluster '%s' created successfully\n", result.Metadata.Name)
	return printSessionCluster(result)
}

func runSessionClusterUpdate(cmd *cobra.Command, args []string) error {
	name := args[0]
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = cfg.Default.Namespace
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	filename, _ := cmd.Flags().GetString("file")
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var sessionCluster api.SessionCluster
	if err := yaml.Unmarshal(data, &sessionCluster); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Set namespace from flag if not in file
	if sessionCluster.Metadata.Namespace == "" {
		sessionCluster.Metadata.Namespace = namespace
	}

	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	result, err := client.UpdateSessionCluster(namespace, name, &sessionCluster)
	if err != nil {
		return err
	}

	fmt.Printf("Session cluster '%s' updated successfully\n", result.Metadata.Name)
	return printSessionCluster(result)
}

func runSessionClusterDelete(cmd *cobra.Command, args []string) error {
	name := args[0]
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

	if err := client.DeleteSessionCluster(namespace, name); err != nil {
		return err
	}

	fmt.Printf("Session cluster '%s' deleted successfully\n", name)
	return nil
}

// Helper functions for printing session clusters
func printSessionClusters(sessionClusters []api.SessionCluster) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(sessionClusters, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(sessionClusters)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tNAMESPACE\tSTATE\tTASKMANAGERS\tFLINK VERSION")
		for _, sc := range sessionClusters {
			state := sc.Status.State
			if state == "" {
				state = sc.Spec.State
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				sc.Metadata.Name,
				sc.Metadata.Namespace,
				state,
				sc.Spec.NumberOfTaskManagers,
				sc.Spec.FlinkVersion,
			)
		}
		w.Flush()
	}
	return nil
}

func printSessionCluster(sc *api.SessionCluster) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(sc, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(sc)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format with details
		fmt.Printf("Name: %s\n", sc.Metadata.Name)
		fmt.Printf("Namespace: %s\n", sc.Metadata.Namespace)
		fmt.Printf("State: %s\n", sc.Status.State)
		fmt.Printf("Desired State: %s\n", sc.Spec.State)
		fmt.Printf("Flink Version: %s\n", sc.Spec.FlinkVersion)
		fmt.Printf("Deployment Target: %s\n", sc.Spec.DeploymentTargetName)
		fmt.Printf("Task Managers: %d\n", sc.Spec.NumberOfTaskManagers)

		if len(sc.Spec.Resources) > 0 {
			fmt.Println("\nResources:")
			for role, res := range sc.Spec.Resources {
				fmt.Printf("  %s:\n", role)
				fmt.Printf("    CPU: %v\n", res.CPU)
				fmt.Printf("    Memory: %v\n", res.Memory)
			}
		}

		if len(sc.Metadata.Labels) > 0 {
			fmt.Println("\nLabels:")
			for k, v := range sc.Metadata.Labels {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}
	}
	return nil
}
