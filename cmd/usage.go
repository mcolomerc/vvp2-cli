package cmd

import (
    "fmt"
    "time"
    "mcolomerc/vvp2cli/pkg/api"
    "github.com/spf13/cobra"
)

var from string
var to string

var usageCmd = &cobra.Command{
    Use:   "usage",
    Short: "Resource usage operations",
}

var usageReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Get resource usage report (platform-wide)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(GetConfig())
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}
		// If from/to not set, default to last 7 days
		f := from
		t := to
		if f == "" || t == "" {
			now := time.Now().UTC()
			layout := "2006-01-02"
			if t == "" {
				t = now.Format(layout)
			}
			if f == "" {
				f = now.Add(-7 * 24 * time.Hour).Format(layout)
			}
		}
		report, err := client.GetResourceUsageReport(f, t)
		if err != nil {
			return fmt.Errorf("failed to get resource usage report: %w", err)
		}
		format := GetConfig().GetOutputFormat()
		switch format {
		case "json":
			// Parse CSV and convert to JSON
			data, err := report.ParseCSV()
			if err != nil {
				return fmt.Errorf("failed to parse CSV data: %w", err)
			}
			return printJSON(data)
		case "yaml":
			// Parse CSV and convert to YAML
			data, err := report.ParseCSV()
			if err != nil {
				return fmt.Errorf("failed to parse CSV data: %w", err)
			}
			return printYAML(data)
		default:
			// For table/default, just print the CSV directly
			fmt.Println(report.CSVData)
			return nil
		}
	},
}

func init() {
	usageReportCmd.Flags().StringVar(&from, "from", "", "Start date (YYYY-MM-DD, inclusive)")
	usageReportCmd.Flags().StringVar(&to, "to", "", "End date (YYYY-MM-DD, exclusive)")
	usageCmd.AddCommand(usageReportCmd)
}
