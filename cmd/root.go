// cmd/root.go
package cmd

import (
    "github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "cribl-storage-tool",
    Short: "A CLI tool to manage Cribl storage",
    Long:  `Cribl Storage Tool is a CLI application to manage various Cribl storage resources.`,
}

// NewRootCmd creates and returns the root command
func NewRootCmd() *cobra.Command {
    return rootCmd
}
