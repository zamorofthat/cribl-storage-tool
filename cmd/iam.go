// cmd/iam.go
package cmd

import (
    "fmt"
    "log"

    // "github.com/aws/aws-sdk-go-v2/config"
    "github.com/spf13/cobra"

    "github.com/zamorofthat/cribl-storage-tool/internal/aws/iam"
    "github.com/zamorofthat/cribl-storage-tool/internal/utils"
)

// iamCmd represents the iam command
var iamCmd = &cobra.Command{
    Use:   "iam",
    Short: "Manage IAM resources",
    Long:  `A subcommand to handle operations related to AWS IAM.`,
}

func init() {
    rootCmd.AddCommand(iamCmd)
    iamCmd.AddCommand(iamSetupCmd)
}

// iamSetupCmd represents the setup command
var iamSetupCmd = &cobra.Command{
    Use:   "setup",
    Short: "Setup IAM role for cross-account access",
    Long:  `A subcommand to setup IAM roles with trust relationships and necessary policies.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Retrieve flags
        roleName, err := cmd.Flags().GetString("role")
        if err != nil {
            log.Fatalf("Error retrieving role flag: %v", err)
        }
        trustedAccountID, err := cmd.Flags().GetString("account")
        if err != nil {
            log.Fatalf("Error retrieving account flag: %v", err)
        }
        externalID, err := cmd.Flags().GetString("external-id")
        if err != nil {
            log.Fatalf("Error retrieving external-id flag: %v", err)
        }

        // Load AWS configuration
        profile, err := cmd.Flags().GetString("profile")
        if err != nil {
            log.Fatalf("Error retrieving profile flag: %v", err)
        }
        region, err := cmd.Flags().GetString("region")
        if err != nil {
            log.Fatalf("Error retrieving region flag: %v", err)
        }

        cfg, err := utils.LoadAWSConfig(profile, region)
        if err != nil {
            log.Fatalf("Unable to load AWS SDK config, %v", err)
        }

        // Initialize IAM client
        iamClient := iam.NewIAMClient(cfg)

        // Setup Trust Relationship
        err = iamClient.SetupTrustRelationship(roleName, trustedAccountID, externalID)
        if err != nil {
            log.Fatalf("Error setting up IAM trust relationship: %v", err)
        }

        fmt.Println("\nIAM trust relationship setup completed successfully.")
    },
}

func init() {
    // Define flags specific to the setup command
    iamSetupCmd.Flags().StringP("role", "r", "CrossAccountAccessRole", "Name of the IAM role to create or update")
    iamSetupCmd.Flags().StringP("account", "a", "", "AWS Account ID to trust (required)")
    iamSetupCmd.Flags().StringP("external-id", "e", "", "External ID for the trust relationship (optional)")
    iamSetupCmd.MarkFlagRequired("account")
}
