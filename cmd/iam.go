package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/zamorofthat/cribl-storage-tool/internal/aws/iam"
	"github.com/zamorofthat/cribl-storage-tool/internal/utils"
)

var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM resources",
	Long:  `A subcommand to handle operations related to AWS IAM.`,
}

func init() {
	rootCmd.AddCommand(iamCmd)
	iamCmd.AddCommand(iamSetupCmd)
}

var iamSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup IAM role for cross-account access",
	Long:  `A subcommand to setup IAM roles with trust relationships and necessary policies.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(os.Stdout).
			With().
			Timestamp().
			Str("command", "iam_setup").
			Logger()

		// Retrieve flags
		roleName, err := cmd.Flags().GetString("role")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving role flag")
		}
		trustedAccountID, err := cmd.Flags().GetString("account")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving account flag")
		}
		externalID, err := cmd.Flags().GetString("external-id")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving external-id flag")
		}
		workspace, err := cmd.Flags().GetString("workspace")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving workspace flag")
		}
		workergroup, err := cmd.Flags().GetString("workergroup")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving workergroup flag")
		}
		action, err := cmd.Flags().GetString("action")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving action flag")
		}
		bucketNames, err := cmd.Flags().GetStringSlice("bucket")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving bucket flag")
		}
		bucketFile, err := cmd.Flags().GetString("bucket-file")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving bucket-file flag")
		}
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving profile flag")
		}
		region, err := cmd.Flags().GetString("region")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving region flag")
		}

		// Handle bucket-file if provided
		if bucketFile != "" {
			data, err := ioutil.ReadFile(bucketFile)
			if err != nil {
				logger.Fatal().Err(err).Str("file", bucketFile).Msg("error reading bucket file")
			}

			var fileBuckets []string
			err = json.Unmarshal(data, &fileBuckets)
			if err != nil {
				logger.Fatal().Err(err).Str("file", bucketFile).Msg("error parsing JSON bucket file")
			}

			bucketNames = append(bucketNames, fileBuckets...)
		}

		// Enforce that at least one of --bucket or --bucket-file is provided
		if len(bucketNames) == 0 {
			logger.Fatal().Msg("at least one bucket name must be provided using the --bucket flag or --bucket-file flag")
		}

		// Load AWS configuration
		cfg, err := utils.LoadAWSConfig(profile, region)
		if err != nil {
			logger.Fatal().Err(err).Msg("unable to load AWS SDK config")
		}

		// Initialize IAM client with logger
		iamClient := iam.NewIAMClient(cfg, logger)

		// Setup Trust Relationship and Policies
		err = iamClient.SetupTrustRelationship(roleName, trustedAccountID, externalID, workspace, workergroup, action, bucketNames)
		if err != nil {
			logger.Fatal().Err(err).Msg("error setting up IAM trust relationship")
		}

		logger.Info().Msg("IAM trust relationship setup completed successfully")
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

	iamSetupCmd.MarkFlagRequired("account")
}
