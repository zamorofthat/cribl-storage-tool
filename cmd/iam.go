// cmd/iam.go
package cmd

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"

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
        workspace, err := cmd.Flags().GetString("workspace")
        if err != nil {
            log.Fatalf("Error retrieving workspace flag: %v", err)
        }
        workergroup, err := cmd.Flags().GetString("workergroup")
        if err != nil {
            log.Fatalf("Error retrieving workergroup flag: %v", err)
        }
        action, err := cmd.Flags().GetString("action")
        if err != nil {
            log.Fatalf("Error retrieving action flag: %v", err)
        }
        bucketNames, err := cmd.Flags().GetStringSlice("bucket")
        if err != nil {
            log.Fatalf("Error retrieving bucket flag: %v", err)
        }
        bucketFile, err := cmd.Flags().GetString("bucket-file")
        if err != nil {
            log.Fatalf("Error retrieving bucket-file flag: %v", err)
        }
        profile, err := cmd.Flags().GetString("profile")
        if err != nil {
            log.Fatalf("Error retrieving profile flag: %v", err)
        }
        region, err := cmd.Flags().GetString("region")
        if err != nil {
            log.Fatalf("Error retrieving region flag: %v", err)
        }

        // Handle bucket-file if provided
        if bucketFile != "" {
            data, err := ioutil.ReadFile(bucketFile)
            if err != nil {
                log.Fatalf("Error reading bucket file: %v", err)
            }

            var fileBuckets []string
            err = json.Unmarshal(data, &fileBuckets)
            if err != nil {
                log.Fatalf("Error parsing JSON bucket file: %v", err)
            }

            bucketNames = append(bucketNames, fileBuckets...)
        }

        // Enforce that at least one of --bucket or --bucket-file is provided
        if len(bucketNames) == 0 {
            log.Fatalf("At least one bucket name must be provided using the --bucket flag or --bucket-file flag.")
        }

        // Load AWS configuration
        cfg, err := utils.LoadAWSConfig(profile, region)
        if err != nil {
            log.Fatalf("Unable to load AWS SDK config, %v", err)
        }

        // Initialize IAM client
        iamClient := iam.NewIAMClient(cfg)

        // Setup Trust Relationship and Policies
        err = iamClient.SetupTrustRelationship(roleName, trustedAccountID, externalID, workspace, workergroup, action, bucketNames)
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
    iamSetupCmd.Flags().StringP("workspace", "w", "main", "Workspace name (default: main)")
    iamSetupCmd.Flags().StringP("workergroup", "g", "default", "Worker group name (default: default)")
    iamSetupCmd.Flags().StringP("action", "s", "search", "Action type for the IAM role (default: search)")
    iamSetupCmd.Flags().StringSliceP("bucket", "b", []string{}, "Name of the S3 bucket to grant access (can specify multiple)")
    iamSetupCmd.Flags().StringP("bucket-file", "f", "", "Path to JSON file containing S3 bucket names (optional)")
    iamSetupCmd.Flags().StringP("profile", "p", "", "AWS profile to use for authentication (optional)")
    iamSetupCmd.Flags().StringP("region", "z", "", "AWS region to target (optional)")

    // Mark required flags
    iamSetupCmd.MarkFlagRequired("account")
    // iamSetupCmd.MarkFlagRequired("bucket") // Removed to allow --bucket-file usage
}