// pkg/aws/s3.go
package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Bucket represents an S3 bucket without the creation date
type Bucket struct {
	Name string `json:"name"`
}

// S3Client wraps the AWS S3 client
type S3Client struct {
	Client *s3.Client
}

// NewS3Client initializes a new S3 client
func NewS3Client(cfg aws.Config) *S3Client {
	return &S3Client{
		Client: s3.NewFromConfig(cfg),
	}
}

// ListBuckets retrieves the list of S3 buckets
func (c *S3Client) ListBuckets() ([]Bucket, error) {
	input := &s3.ListBucketsInput{}
	result, err := c.Client.ListBuckets(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	for _, b := range result.Buckets {
		buckets = append(buckets, Bucket{
			Name: aws.ToString(b.Name),
		})
	}
	return buckets, nil
}

// PrintBucketsText prints the list of buckets in text format
func (c *S3Client) PrintBucketsText(buckets []Bucket) {
	fmt.Println("Listing S3 Buckets:")
	for _, bucket := range buckets {
		fmt.Printf(" - %s\n", bucket.Name)
	}
}

// PrintBucketsJSON prints the list of buckets in JSON format
func (c *S3Client) PrintBucketsJSON(buckets []Bucket) error {
	// Marshal the slice to JSON
	jsonData, err := json.MarshalIndent(buckets, "", "  ")
	if err != nil {
		return err
	}

	// Print the JSON data
	fmt.Println(string(jsonData))
	return nil
}

// Add this new function
// PrintBucketsNameOnly prints just the bucket names, one per line
func (c *S3Client) PrintBucketsNameOnly(buckets []Bucket) {
    for _, bucket := range buckets {
        fmt.Println(bucket.Name)
    }
}