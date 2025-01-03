// cmd/list.go
package cmd

import (
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/spf13/cobra"
    "github.com/zamorofthat/cribl-storage-tool/internal/aws/s3"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all S3 buckets",
    Long:  `A subcommand to list all AWS S3 buckets.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Retrieve flags
        profile, err := cmd.Flags().GetString("profile")
        if err != nil {
            log.Fatalf("Error retrieving profile flag: %v", err)
        }
        region, err := cmd.Flags().GetString("region")
        if err != nil {
            log.Fatalf("Error retrieving region flag: %v", err)
        }

        // Load AWS configuration
        cfg, err := loadAWSConfig(profile, region)
        if err != nil {
            log.Fatalf("Unable to load AWS SDK config, %v", err)
        }

        // Initialize S3 client
        s3Client := s3.NewS3Client(cfg)

        // List S3 Buckets
        fmt.Println("\nListing S3 Buckets:")
        buckets, err := s3Client.ListBuckets()
        if err != nil {
            log.Fatalf("Error listing S3 buckets: %v", err)
        }
        for _, bucket := range buckets {
            fmt.Println(" -", bucket)
        }
    },
}

// loadAWSConfig loads the AWS configuration with optional profile and region
func loadAWSConfig(profile, region string) (aws.Config, error) {
    var cfg aws.Config
    var err error

    options := []func(*config.LoadOptions) error{}

    if profile != "" {
        options = append(options, config.WithSharedConfigProfile(profile))
    }
    if region != "" {
        options = append(options, config.WithRegion(region))
    }

    cfg, err = config.LoadDefaultConfig(context.TODO(), options...)
    return cfg, err
}

func init() {
    s3Cmd.AddCommand(listCmd)

    // Define flags specific to the list command
    listCmd.Flags().StringP("profile", "p", "", "AWS profile to use for authentication")
    listCmd.Flags().StringP("region", "r", "", "AWS region to target")

    // Optionally, mark flags as required
    // listCmd.MarkFlagRequired("profile")
    // listCmd.MarkFlagRequired("region")
}
