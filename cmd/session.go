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
	sessionNamespace string
	sessionFile      string
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:     "session",
	Aliases: []string{"sessions"},
	Short:   "Manage VVP sessions",
	Long:    `Create, read, update, and delete Ververica Platform sessions (SQL sessions).`,
}

// listSessionsCmd lists all sessions
var listSessionsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all sessions",
	RunE:    runListSessions,
}

// getSessionCmd gets a session by name
var getSessionCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a session by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetSession,
}

// createSessionCmd creates a new session
var createSessionCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new session from a YAML/JSON file",
	RunE:  runCreateSession,
}

// updateSessionCmd updates a session
var updateSessionCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update an existing session",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdateSession,
}

// deleteSessionCmd deletes a session
var deleteSessionCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a session",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeleteSession,
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(listSessionsCmd)
	sessionCmd.AddCommand(getSessionCmd)
	sessionCmd.AddCommand(createSessionCmd)
	sessionCmd.AddCommand(updateSessionCmd)
	sessionCmd.AddCommand(deleteSessionCmd)

	// Flags for session commands
	sessionCmd.PersistentFlags().StringVarP(&sessionNamespace, "namespace", "n", "", "Namespace (defaults to config if not set)")

	createSessionCmd.Flags().StringVarP(&sessionFile, "file", "f", "", "Path to session YAML/JSON file (required)")
	createSessionCmd.MarkFlagRequired("file")

	updateSessionCmd.Flags().StringVarP(&sessionFile, "file", "f", "", "Path to session YAML/JSON file (required)")
	updateSessionCmd.MarkFlagRequired("file")
}

func runListSessions(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveSessionNamespace()
	if err != nil {
		return err
	}

	sessions, err := client.ListSessions(ns)
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	return printSessions(sessions.Items)
}

func runGetSession(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveSessionNamespace()
	if err != nil {
		return err
	}

	session, err := client.GetSession(ns, args[0])
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	return printSession(session)
}

func runCreateSession(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveSessionNamespace()
	if err != nil {
		return err
	}

	session, err := loadSessionFromFile(sessionFile)
	if err != nil {
		return err
	}

	result, err := client.CreateSession(ns, session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	fmt.Printf("Session %s created successfully\n", result.Metadata.Name)
	return printSession(result)
}

func runUpdateSession(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveSessionNamespace()
	if err != nil {
		return err
	}

	session, err := loadSessionFromFile(sessionFile)
	if err != nil {
		return err
	}

	result, err := client.UpdateSession(ns, args[0], session)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	fmt.Printf("Session %s updated successfully\n", result.Metadata.Name)
	return printSession(result)
}

func runDeleteSession(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	ns, err := effectiveSessionNamespace()
	if err != nil {
		return err
	}

	if err := client.DeleteSession(ns, args[0]); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	fmt.Printf("Session %s deleted successfully\n", args[0])
	return nil
}

func loadSessionFromFile(filename string) (*api.Session, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var session api.Session

	// Try JSON first
	if err := json.Unmarshal(data, &session); err != nil {
		// Try YAML
		if err := yaml.Unmarshal(data, &session); err != nil {
			return nil, fmt.Errorf("failed to parse file as JSON or YAML: %w", err)
		}
	}

	return &session, nil
}

func printSessions(sessions []api.Session) error {
	if len(sessions) == 0 {
		fmt.Println("No sessions found")
		return nil
	}

	format := GetConfig().GetOutputFormat()

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(sessions)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(sessions)
	default:
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tNAMESPACE\tSTATE\tCREATED")
		for _, s := range sessions {
			state := s.Status.State
			if state == "" {
				state = "N/A"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				s.Metadata.Name,
				s.Metadata.Namespace,
				state,
				s.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		return w.Flush()
	}
}

func printSession(session *api.Session) error {
	format := GetConfig().GetOutputFormat()

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(session)
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(session)
	default:
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(2)
		return encoder.Encode(session)
	}
}

// effectiveSessionNamespace returns the namespace for session commands,
// preferring the --namespace flag, then falling back to the default from config.
// If neither is set, it returns a helpful error.
func effectiveSessionNamespace() (string, error) {
	if sessionNamespace != "" {
		return sessionNamespace, nil
	}
	if ns := GetConfig().GetNamespace(); ns != "" {
		return ns, nil
	}
	return "", fmt.Errorf("namespace not specified. Provide --namespace or set default.namespace in ~/.vvp2/config.yaml (or VVP_DEFAULT_NAMESPACE)")
}
