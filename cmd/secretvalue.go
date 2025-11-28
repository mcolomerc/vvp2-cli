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

var secretValueCmd = &cobra.Command{
	Use:     "secret-value",
	Aliases: []string{"secret", "secrets", "sv"},
	Short:   "Manage VVP secret values",
	Long:    `Create, list, view, update, and delete Ververica Platform secret values. Secret values store sensitive data like passwords and API keys.`,
}

var secretValueListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List secret values in a namespace",
	RunE:    runSecretValueList,
}

var secretValueGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a secret value by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runSecretValueGet,
}

var secretValueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a secret value from a file",
	RunE:  runSecretValueCreate,
}

var secretValueUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update a secret value",
	Args:  cobra.ExactArgs(1),
	RunE:  runSecretValueUpdate,
}

var secretValueDeleteCmd = &cobra.Command{
	Use:     "delete [name]",
	Aliases: []string{"rm"},
	Short:   "Delete a secret value",
	Args:    cobra.ExactArgs(1),
	RunE:    runSecretValueDelete,
}

func init() {
	rootCmd.AddCommand(secretValueCmd)

	secretValueCmd.AddCommand(secretValueListCmd)
	secretValueCmd.AddCommand(secretValueGetCmd)
	secretValueCmd.AddCommand(secretValueCreateCmd)
	secretValueCmd.AddCommand(secretValueUpdateCmd)
	secretValueCmd.AddCommand(secretValueDeleteCmd)

	// Add flags
	secretValueListCmd.Flags().StringP("namespace", "n", "", "Namespace")
	secretValueGetCmd.Flags().StringP("namespace", "n", "", "Namespace")

	secretValueCreateCmd.Flags().StringP("namespace", "n", "", "Namespace")
	secretValueCreateCmd.Flags().StringP("file", "f", "", "File containing secret value definition")
	secretValueCreateCmd.MarkFlagRequired("file")

	secretValueUpdateCmd.Flags().StringP("namespace", "n", "", "Namespace")
	secretValueUpdateCmd.Flags().StringP("file", "f", "", "File containing secret value definition")
	secretValueUpdateCmd.MarkFlagRequired("file")

	secretValueDeleteCmd.Flags().StringP("namespace", "n", "", "Namespace")
}

func runSecretValueList(cmd *cobra.Command, args []string) error {
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

	secretValues, err := client.ListSecretValues(namespace)
	if err != nil {
		return err
	}

	return printSecretValues(secretValues.Items)
}

func runSecretValueGet(cmd *cobra.Command, args []string) error {
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

	secretValue, err := client.GetSecretValue(namespace, name)
	if err != nil {
		return err
	}

	return printSecretValue(secretValue)
}

func runSecretValueCreate(cmd *cobra.Command, args []string) error {
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

	var secretValue api.SecretValue
	if err := yaml.Unmarshal(data, &secretValue); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Set namespace from flag if not in file
	if secretValue.Metadata.Namespace == "" {
		secretValue.Metadata.Namespace = namespace
	}

	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	result, err := client.CreateSecretValue(namespace, &secretValue)
	if err != nil {
		return err
	}

	fmt.Printf("Secret value '%s' created successfully\n", result.Metadata.Name)
	return printSecretValue(result)
}

func runSecretValueUpdate(cmd *cobra.Command, args []string) error {
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

	var secretValue api.SecretValue
	if err := yaml.Unmarshal(data, &secretValue); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Set namespace from flag if not in file
	if secretValue.Metadata.Namespace == "" {
		secretValue.Metadata.Namespace = namespace
	}

	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	result, err := client.UpdateSecretValue(namespace, name, &secretValue)
	if err != nil {
		return err
	}

	fmt.Printf("Secret value '%s' updated successfully\n", result.Metadata.Name)
	return printSecretValue(result)
}

func runSecretValueDelete(cmd *cobra.Command, args []string) error {
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

	if err := client.DeleteSecretValue(namespace, name); err != nil {
		return err
	}

	fmt.Printf("Secret value '%s' deleted successfully\n", name)
	return nil
}

// Helper functions for printing secret values
func printSecretValues(secretValues []api.SecretValue) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(secretValues, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(secretValues)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format
		if len(secretValues) == 0 {
			fmt.Println("No secret values found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tNAMESPACE\tKIND\tCREATED")
		for _, sv := range secretValues {
			kind := sv.Spec.Kind
			if kind == "" {
				kind = "-"
			}

			created := "-"
			if !sv.Metadata.CreatedAt.IsZero() {
				created = sv.Metadata.CreatedAt.Format("2006-01-02 15:04:05")
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				sv.Metadata.Name,
				sv.Metadata.Namespace,
				kind,
				created,
			)
		}
		w.Flush()
	}
	return nil
}

func printSecretValue(sv *api.SecretValue) error {
	outputFormat, _ := rootCmd.PersistentFlags().GetString("output")

	switch outputFormat {
	case "json":
		data, err := json.MarshalIndent(sv, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(sv)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		// Table format with details (but don't print the actual secret value!)
		fmt.Printf("Name: %s\n", sv.Metadata.Name)
		if sv.Metadata.ID != "" {
			fmt.Printf("ID: %s\n", sv.Metadata.ID)
		}
		fmt.Printf("Namespace: %s\n", sv.Metadata.Namespace)

		if sv.Spec.Kind != "" {
			fmt.Printf("Kind: %s\n", sv.Spec.Kind)
		}

		// Don't print the actual secret value in default output
		if sv.Spec.Value != "" {
			fmt.Printf("Value: <hidden> (use -o json or -o yaml to view)\n")
		}

		if len(sv.Metadata.Labels) > 0 {
			fmt.Println("\nLabels:")
			for k, v := range sv.Metadata.Labels {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}

		if len(sv.Metadata.Annotations) > 0 {
			fmt.Println("\nAnnotations:")
			for k, v := range sv.Metadata.Annotations {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}

		if !sv.Metadata.CreatedAt.IsZero() {
			fmt.Printf("\nCreated At: %s\n", sv.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		if !sv.Metadata.ModifiedAt.IsZero() {
			fmt.Printf("Modified At: %s\n", sv.Metadata.ModifiedAt.Format("2006-01-02 15:04:05"))
		}
	}
	return nil
}
