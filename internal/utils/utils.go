// internal/utils/aws.go
package utils

import (
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
)

// LoadAWSConfig loads the AWS configuration with optional profile and region
func LoadAWSConfig(profile, region string) (aws.Config, error) {
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
