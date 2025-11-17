package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display version, commit hash, and build time information for vvp2 CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vvp2 CLI\n")
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Commit:     %s\n", commit)
		fmt.Printf("Built:      %s\n", buildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
