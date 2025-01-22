// cmd/s3.go
package cmd

import (
	"github.com/spf13/cobra"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Manage S3 resources",
	Long:  `A subcommand to handle operations related to AWS S3.`,
}

func init() {
	// Add the s3 command to the root command
	rootCmd.AddCommand(s3Cmd)

	// Add the list subcommand to the s3 command
	s3Cmd.AddCommand(listCmd)
}