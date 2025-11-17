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
	namespaceFile string
)

// namespaceCmd represents the namespace command
var namespaceCmd = &cobra.Command{
	Use:     "namespace",
	Aliases: []string{"ns", "namespaces"},
	Short:   "Manage VVP namespaces",
	Long:    `Create, read, update, and delete Ververica Platform namespaces.`,
}

// listNamespacesCmd lists all namespaces
var listNamespacesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all namespaces",
	RunE:    runListNamespaces,
}

// getNamespaceCmd gets a namespace by name
var getNamespaceCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a namespace by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetNamespace,
}

// createNamespaceCmd creates a new namespace
var createNamespaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new namespace from a YAML/JSON file",
	RunE:  runCreateNamespace,
}

// updateNamespaceCmd updates a namespace
var updateNamespaceCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update an existing namespace",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdateNamespace,
}

// deleteNamespaceCmd deletes a namespace
var deleteNamespaceCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a namespace",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeleteNamespace,
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
	namespaceCmd.AddCommand(listNamespacesCmd)
	namespaceCmd.AddCommand(getNamespaceCmd)
	namespaceCmd.AddCommand(createNamespaceCmd)
	namespaceCmd.AddCommand(updateNamespaceCmd)
	namespaceCmd.AddCommand(deleteNamespaceCmd)

	// Flags for namespace commands
	createNamespaceCmd.Flags().StringVarP(&namespaceFile, "file", "f", "", "Path to namespace YAML/JSON file (required)")
	createNamespaceCmd.MarkFlagRequired("file")

	updateNamespaceCmd.Flags().StringVarP(&namespaceFile, "file", "f", "", "Path to namespace YAML/JSON file (required)")
	updateNamespaceCmd.MarkFlagRequired("file")
}

func runListNamespaces(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	namespaces, err := client.ListNamespaces()
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	return printNamespaces(namespaces.Items)
}

func runGetNamespace(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	namespace, err := client.GetNamespace(args[0])
	if err != nil {
		return fmt.Errorf("failed to get namespace: %w", err)
	}

	return printNamespace(namespace)
}

func runCreateNamespace(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	namespace, err := loadNamespaceFromFile(namespaceFile)
	if err != nil {
		return err
	}

	result, err := client.CreateNamespace(namespace)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	fmt.Printf("Namespace %s created successfully\n", result.Metadata.Name)
	return printNamespace(result)
}

func runUpdateNamespace(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	namespace, err := loadNamespaceFromFile(namespaceFile)
	if err != nil {
		return err
	}

	result, err := client.UpdateNamespace(args[0], namespace)
	if err != nil {
		return fmt.Errorf("failed to update namespace: %w", err)
	}

	fmt.Printf("Namespace %s updated successfully\n", result.Metadata.Name)
	return printNamespace(result)
}

func runDeleteNamespace(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	if err := client.DeleteNamespace(args[0]); err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	fmt.Printf("Namespace %s deleted successfully\n", args[0])
	return nil
}

func loadNamespaceFromFile(filename string) (*api.Namespace, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var namespace api.Namespace

	// Try JSON first
	if err := json.Unmarshal(data, &namespace); err != nil {
		// Try YAML
		if err := yaml.Unmarshal(data, &namespace); err != nil {
			return nil, fmt.Errorf("failed to parse file as JSON or YAML: %w", err)
		}
	}

	return &namespace, nil
}

func printNamespaces(namespaces []api.Namespace) error {
	if len(namespaces) == 0 {
		fmt.Println("No namespaces found")
		return nil
	}

	format := GetConfig().GetOutputFormat()

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(namespaces)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(namespaces)
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATE\tCREATED")
		for _, ns := range namespaces {
			state := ns.Status.State
			if state == "" {
				state = "N/A"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				ns.Metadata.Name,
				state,
				ns.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		return w.Flush()
	}
}

func printNamespace(namespace *api.Namespace) error {
	format := GetConfig().GetOutputFormat()

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(namespace)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(namespace)
	default:
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(namespace)
	}
}
