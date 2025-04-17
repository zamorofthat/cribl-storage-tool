package cmd

import (
	"encoding/json"
	"fmt"

	//     "io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	criblawshelper "github.com/zamorofthat/cribl-storage-tool/pkg/aws"
	"github.com/zamorofthat/cribl-storage-tool/pkg/utils"
	// "github.com/zamorofthat/cribl-storage-tool/internal/utils"
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

func parseWorkerArn(arn string) (accountID, workspace, workergroup string, err error) {
	// Expected format: arn:aws:iam::ACCOUNT:role/WORKSPACE-WORKERGROUP
	parts := strings.Split(arn, ":")
	if len(parts) != 6 {
		return "", "", "", fmt.Errorf("invalid ARN format")
	}

	// Extract account ID
	accountID = parts[4]

	// Extract workspace and workergroup
	roleParts := strings.Split(parts[5], "/")
	if len(roleParts) != 2 {
		return "", "", "", fmt.Errorf("invalid role format in ARN")
	}

	// Split the role name into workspace and workergroup
	nameComponents := strings.Split(roleParts[1], "-")
	if len(nameComponents) != 2 {
		return "", "", "", fmt.Errorf("invalid role name format, expected workspace-workergroup")
	}

	workspace = nameComponents[0]
	workergroup = nameComponents[1]

	return accountID, workspace, workergroup, nil
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

		// Get the worker ARN if provided
		workerArn, err := cmd.Flags().GetString("cribl-worker-arn")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving cribl-worker-arn flag")
		}

		var trustedAccountID, workspace, workergroup string

		if workerArn != "" {
			// Parse the worker ARN
			trustedAccountID, workspace, workergroup, err = parseWorkerArn(workerArn)
			if err != nil {
				logger.Fatal().Err(err).Str("arn", workerArn).Msg("failed to parse worker ARN")
			}
			logger.Info().
				Str("account_id", trustedAccountID).
				Str("workspace", workspace).
				Str("workergroup", workergroup).
				Msg("parsed worker ARN")
		} else {
			// Use individual flags if ARN not provided
			trustedAccountID, err = cmd.Flags().GetString("account")
			if err != nil {
				logger.Fatal().Err(err).Msg("error retrieving account flag")
			}
			workspace, err = cmd.Flags().GetString("workspace")
			if err != nil {
				logger.Fatal().Err(err).Msg("error retrieving workspace flag")
			}
			workergroup, err = cmd.Flags().GetString("workergroup")
			if err != nil {
				logger.Fatal().Err(err).Msg("error retrieving workergroup flag")
			}
		}

		roleName, err := cmd.Flags().GetString("role")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving role flag")
		}

		externalID, err := cmd.Flags().GetString("external-id")
		if err != nil {
			logger.Fatal().Err(err).Msg("error retrieving external-id flag")
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
			data, err := os.ReadFile(bucketFile)
			if err != nil {
				logger.Fatal().Err(err).Str("file", bucketFile).Msg("error reading bucket file")
			}

			// Try to parse as JSON first
			var fileBuckets []string
			jsonErr := json.Unmarshal(data, &fileBuckets)
			if jsonErr == nil {
				// Successfully parsed as simple JSON array
				bucketNames = append(bucketNames, fileBuckets...)
				logger.Info().Strs("buckets", fileBuckets).Msg("loaded buckets from JSON array file")
			} else {
				// Try to parse as array of objects with name field
				var bucketObjects []struct {
					Name string `json:"name"`
				}
				jsonErr = json.Unmarshal(data, &bucketObjects)
				if jsonErr == nil {
					// Successfully parsed as array of objects
					for _, obj := range bucketObjects {
						bucketNames = append(bucketNames, obj.Name)
					}
					logger.Info().Strs("buckets", bucketNames).Msg("loaded buckets from JSON objects file")
				} else {
					// If both JSON parsing attempts failed, try to parse as plain text (one bucket per line)
					lines := strings.Split(string(data), "\n")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if line != "" {
							bucketNames = append(bucketNames, line)
						}
					}
					logger.Info().Strs("buckets", bucketNames).Msg("loaded buckets from text file")
				}
			}
		}

		// Enforce that at least one of --bucket or --bucket-file is provided
		if len(bucketNames) == 0 {
			logger.Fatal().Msg("at least one bucket name must be provided using the --bucket flag or --bucket-file flag")
		}

		// Load AWS configuration
		cfg, err := utils.LoadAWSConfig(cmd.Context(), profile, region, logger)
		if err != nil {
			logger.Fatal().Err(err).
				Str("profile", profile).
				Str("region", region).
				Msg("unable to load AWS SDK config")
		}

		// Initialize IAM client with logger
		iamClient := criblawshelper.NewIAMClient(cfg, logger)

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
	iamSetupCmd.Flags().String("cribl-worker-arn", "", "Cribl worker ARN (e.g., arn:aws:iam::ACCOUNT:role/WORKSPACE-WORKERGROUP)")
	iamSetupCmd.Flags().StringP("role", "r", "CrossAccountAccessRole", "Name of the IAM role to create or update")
	iamSetupCmd.Flags().StringP("account", "a", "", "AWS Account ID to trust (required if --cribl-worker-arn not provided)")
	iamSetupCmd.Flags().StringP("external-id", "e", "", "External ID for the trust relationship (optional)")
	iamSetupCmd.Flags().StringP("workspace", "w", "main", "Workspace name (default: main)")
	iamSetupCmd.Flags().StringP("workergroup", "g", "default", "Worker group name (default: default)")
	iamSetupCmd.Flags().StringP("action", "s", "search", "Action type for the IAM role (default: search)")
	iamSetupCmd.Flags().StringSliceP("bucket", "b", []string{}, "Name of the S3 bucket to grant access (can specify multiple)")
	iamSetupCmd.Flags().StringP("bucket-file", "f", "", "Path to JSON file containing S3 bucket names (optional)")
	iamSetupCmd.Flags().StringP("profile", "p", "", "AWS profile to use for authentication (optional)")
	iamSetupCmd.Flags().StringP("region", "z", "", "AWS region to target (optional)")

	// Only require account if cribl-worker-arn is not provided
	// iamSetupCmd.MarkFlagRequired("account")
}
