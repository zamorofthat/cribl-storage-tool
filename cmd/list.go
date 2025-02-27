// cmd/list.go
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	criblawshelper "github.com/zamorofthat/cribl-storage-tool/pkg/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"context" // Ensure this line is present
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all S3 buckets",
	Long:  `A subcommand to list all AWS S3 buckets.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve flags
		outputFormat, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Fatalf("Error retrieving output flag: %v", err)
		}
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			log.Fatalf("Error retrieving profile flag: %v", err)
		}
		region, err := cmd.Flags().GetString("region")
		if err != nil {
			log.Fatalf("Error retrieving region flag: %v", err)
		}
		filter, err := cmd.Flags().GetString("filter")
		if err != nil {
			log.Fatalf("Error retrieving filter flag: %v", err)
		}
		regexPattern, err := cmd.Flags().GetString("regex")
		if err != nil {
			log.Fatalf("Error retrieving regex flag: %v", err)
		}
		bucketFile, err := cmd.Flags().GetString("bucket-file")
		if err != nil {
			log.Fatalf("Error retrieving bucket-file flag: %v", err)
		}

		// Enforce mutual exclusivity between --filter, --regex, and --bucket-file
		count := 0
		if filter != "" {
			count++
		}
		if regexPattern != "" {
			count++
		}
		if bucketFile != "" {
			count++
		}
		if count > 1 {
			log.Fatalf("Flags --filter, --regex, and --bucket-file cannot be used together. Please use only one.")
		}

		// Load AWS configuration
		cfg, err := loadAWSConfig(profile, region)
		if err != nil {
			log.Fatalf("Unable to load AWS SDK config: %v", err)
		}

		// Initialize S3 client
		s3Client := criblawshelper.NewS3Client(cfg)

		// Retrieve the list of buckets
		buckets, err := s3Client.ListBuckets()
		if err != nil {
			log.Fatalf("Error listing S3 buckets: %v", err)
		}

		// Handle --bucket-file if provided
		if bucketFile != "" {
			fileBuckets, err := loadBucketsFromFile(bucketFile)
			if err != nil {
				log.Fatalf("Error loading buckets from file: %v", err)
			}
			// Override buckets with those from the file
			buckets = fileBuckets
		}

		// Apply substring filter if provided
		if filter != "" {
			var filteredBuckets []criblawshelper.Bucket
			for _, bucket := range buckets {
				if strings.Contains(bucket.Name, filter) {
					filteredBuckets = append(filteredBuckets, bucket)
				}
			}
			buckets = filteredBuckets
		}

		// Apply regex filter if provided
		if regexPattern != "" {
			compiledRegex, err := regexp.Compile(regexPattern)
			if err != nil {
				log.Fatalf("Invalid regex pattern: %v", err)
			}
			var regexFilteredBuckets []criblawshelper.Bucket
			for _, bucket := range buckets {
				if compiledRegex.MatchString(bucket.Name) {
					regexFilteredBuckets = append(regexFilteredBuckets, bucket)
				}
			}
			buckets = regexFilteredBuckets
		}

		// Format and print the output
        switch outputFormat {
        case "json":
            err = s3Client.PrintBucketsJSON(buckets)
            if err != nil {
                log.Fatalf("Error printing buckets in JSON format: %v", err)
            }
        case "names":
            s3Client.PrintBucketsNameOnly(buckets)
        case "text":
            fallthrough
        default:
            s3Client.PrintBucketsText(buckets)
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

// loadBucketsFromFile reads a JSON file containing a list of bucket names
func loadBucketsFromFile(filePath string) ([]criblawshelper.Bucket, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var bucketNames []string
	err = json.Unmarshal(data, &bucketNames)
	if err != nil {
		return nil, err
	}

	var buckets []criblawshelper.Bucket
	for _, name := range bucketNames {
		buckets = append(buckets, criblawshelper.Bucket{Name: name})
	}
	return buckets, nil
}

func init() {
	// Define flags specific to the list command
    listCmd.Flags().StringP("output", "o", "text", "Output format: text, json, or names")
	listCmd.Flags().StringP("profile", "p", "", "AWS profile to use for authentication (optional)")
	listCmd.Flags().StringP("region", "r", "", "AWS region to target (optional)")
	listCmd.Flags().StringP("filter", "f", "", "Filter bucket names containing the specified substring (optional)")
	listCmd.Flags().StringP("regex", "x", "", "Filter bucket names matching the specified regular expression (optional)")
	listCmd.Flags().StringP("bucket-file", "b", "", "Path to file containing S3 bucket names (optional)")
}