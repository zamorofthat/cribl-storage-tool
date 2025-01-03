// internal/aws/s3/s3.go
package s3

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

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

// ListBuckets lists all S3 buckets
func (c *S3Client) ListBuckets() ([]string, error) {
    output, err := c.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
    if err != nil {
        return nil, fmt.Errorf("unable to list S3 buckets: %w", err)
    }

    var buckets []string
    for _, bucket := range output.Buckets {
        buckets = append(buckets, *bucket.Name)
    }
    return buckets, nil
}
