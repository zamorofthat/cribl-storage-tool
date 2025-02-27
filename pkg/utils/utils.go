// pkg/utils/utils.go
package utils

import (
    "context"
    "github.com/rs/zerolog"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
)

// LoadAWSConfig loads the AWS configuration with optional profile and region
// internal/utils/aws.go
// internal/utils/aws.go
func LoadAWSConfig(ctx context.Context, profile, region string, logger zerolog.Logger) (aws.Config, error) {
    var cfg aws.Config
    var err error

    options := []func(*config.LoadOptions) error{}

    if profile != "" {
        logger.Debug().Str("profile", profile).Msg("using specified AWS profile")
        options = append(options, config.WithSharedConfigProfile(profile))
    } else {
        logger.Debug().Msg("using default AWS profile")
    }

    if region != "" {
        logger.Debug().Str("region", region).Msg("using specified AWS region")
        options = append(options, config.WithRegion(region))
    } else {
        logger.Debug().Msg("region not specified, will use from AWS profile or environment")
    }

    cfg, err = config.LoadDefaultConfig(ctx, options...)
    if err != nil {
        return cfg, err
    }

    // Log the region that was actually resolved
    if region == "" {
        logger.Info().Str("resolved_region", cfg.Region).Msg("using region from AWS profile or environment")
    } else {
        logger.Info().Str("region", cfg.Region).Msg("AWS config loaded with specified region")
    }

    return cfg, nil
}