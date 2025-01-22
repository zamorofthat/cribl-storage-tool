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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	// Here you can define persistent flags and configuration settings.
	// For example, a persistent flag for verbose logging can be added here.
}