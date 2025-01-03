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
    rootCmd.AddCommand(s3Cmd)
}
