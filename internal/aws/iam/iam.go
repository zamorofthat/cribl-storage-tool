package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/rs/zerolog"
)

type IAMClient struct {
	Client *iam.Client
	logger zerolog.Logger
}

type TrustPolicyDocument struct {
	Version   string           `json:"Version"`
	Statement []TrustStatement `json:"Statement"`
}

type TrustStatement struct {
	Effect    string    `json:"Effect"`
	Principal Principal `json:"Principal"`
	Action    []string  `json:"Action"`
	Condition Condition `json:"Condition,omitempty"`
}

type RolePolicyDocument struct {
	Version   string          `json:"Version"`
	Statement []RoleStatement `json:"Statement"`
}

type RoleStatement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
}

type Principal struct {
	AWS string `json:"AWS"`
}

type Condition struct {
	StringEquals map[string]string `json:"StringEquals,omitempty"`
}

func NewIAMClient(cfg aws.Config, logger zerolog.Logger) *IAMClient {
	return &IAMClient{
		Client: iam.NewFromConfig(cfg),
		logger: logger.With().Str("component", "iam_client").Logger(),
	}
}

func (c *IAMClient) SetupTrustRelationship(roleName, trustedAccountID, externalID string, workspace string,
	workergroup string, action string, bucketNames []string) error {
	logger := c.logger.With().
		Str("role_name", roleName).
		Str("trusted_account_id", trustedAccountID).
		Str("workspace", workspace).
		Str("workergroup", workergroup).
		Str("action", action).
		Strs("bucket_names", bucketNames).
		Logger()

	logger.Info().Msg("setting up trust relationship")

	if err := c.validateInputs(roleName, trustedAccountID, workspace, workergroup, bucketNames); err != nil {
		logger.Error().Err(err).Msg("input validation failed")
		return err
	}

	trustPolicy := c.createTrustPolicy(trustedAccountID, workspace, workergroup, externalID, action)
	logger.Debug().RawJSON("trust_policy", []byte(trustPolicy)).Msg("created trust policy")

	if err := c.ensureRoleExists(roleName, trustPolicy); err != nil {
		logger.Error().Err(err).Msg("failed to ensure role exists")
		return err
	}

	if err := c.attachS3Policies(roleName, bucketNames); err != nil {
		logger.Error().Err(err).Msg("failed to attach S3 policies")
		return err
	}

	logger.Info().Msg("successfully set up trust relationship")
	return nil
}

func (c *IAMClient) validateInputs(roleName, trustedAccountID, workspace, workergroup string, bucketNames []string) error {
	logger := c.logger.With().
		Str("role_name", roleName).
		Str("trusted_account_id", trustedAccountID).
		Str("workspace", workspace).
		Str("workergroup", workergroup).
		Strs("bucket_names", bucketNames).
		Logger()

	logger.Debug().Msg("validating inputs")

	if roleName == "" {
		logger.Error().Msg("role name is empty")
		return fmt.Errorf("roleName cannot be empty")
	}
	if trustedAccountID == "" {
		logger.Error().Msg("trusted account ID is empty")
		return fmt.Errorf("trustedAccountID cannot be empty")
	}
	if len(bucketNames) == 0 {
		logger.Error().Msg("no bucket names provided")
		return fmt.Errorf("at least one bucketName must be provided")
	}

	logger.Debug().Msg("input validation successful")
	return nil
}

func (c *IAMClient) createTrustPolicy(trustedAccountID, workspace, workergroup, externalID, action string) string {
	logger := c.logger.With().
		Str("trusted_account_id", trustedAccountID).
		Str("workspace", workspace).
		Str("workergroup", workergroup).
		Str("action", action).
		Logger()

	logger.Debug().Msg("creating trust policy")

	policy := TrustPolicyDocument{
		Version: "2012-10-17",
		Statement: []TrustStatement{
			{
				Effect: "Allow",
				Action: []string{"sts:AssumeRole", "sts:TagSession", "sts:SetSourceIdentity"},
				Condition: Condition{
					StringEquals: map[string]string{
						"sts:ExternalId": externalID,
					},
				},
			},
		},
	}

	if action != "search" {
		policy.Statement[0].Principal = Principal{
			AWS: fmt.Sprintf("arn:aws:iam::%s:role/%s-%s", trustedAccountID, workspace, workergroup),
		}
	} else {
		policy.Statement[0].Principal = Principal{
			AWS: fmt.Sprintf("arn:aws:iam::%s:role/search-exec-%s", trustedAccountID, workspace),
		}
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal policy")
		return ""
	}

	logger.Debug().RawJSON("policy", policyJSON).Msg("trust policy created")
	return string(policyJSON)
}

func (c *IAMClient) ensureRoleExists(roleName, trustPolicy string) error {
	logger := c.logger.With().Str("role_name", roleName).Logger()

	_, err := c.Client.GetRole(context.TODO(), &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})

	if err != nil {
		var noSuchEntity *types.NoSuchEntityException
		if errors.As(err, &noSuchEntity) {
			logger.Info().Msg("role does not exist, creating new role")
			return c.createRole(roleName, trustPolicy)
		}
		logger.Error().Err(err).Msg("failed to get IAM role")
		return fmt.Errorf("failed to get IAM role: %w", err)
	}

	logger.Info().Msg("updating existing role trust policy")
	return c.updateRoleTrustPolicy(roleName, trustPolicy)
}

func (c *IAMClient) createRole(roleName, trustPolicy string) error {
	logger := c.logger.With().Str("role_name", roleName).Logger()

	_, err := c.Client.CreateRole(context.TODO(), &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(trustPolicy),
		Description:              aws.String("Role for cross-account access to S3"),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create IAM role")
		return fmt.Errorf("failed to create IAM role: %w", err)
	}

	logger.Info().Msg("created IAM role with trust relationship")
	return nil
}

func (c *IAMClient) updateRoleTrustPolicy(roleName, trustPolicy string) error {
	logger := c.logger.With().Str("role_name", roleName).Logger()

	_, err := c.Client.UpdateAssumeRolePolicy(context.TODO(), &iam.UpdateAssumeRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyDocument: aws.String(trustPolicy),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update trust policy")
		return fmt.Errorf("failed to update trust policy for role '%s': %w", roleName, err)
	}

	logger.Info().Msg("updated trust relationship for IAM role")
	return nil
}

func (c *IAMClient) attachS3Policies(roleName string, bucketNames []string) error {
	logger := c.logger.With().
		Str("role_name", roleName).
		Strs("bucket_names", bucketNames).
		Logger()

	policyName := "CrossAccountAccessPolicy"
	policyDocument := c.createS3PolicyDocument(bucketNames)

	logger.Debug().RawJSON("policy_document", []byte(policyDocument)).Msg("creating S3 policy")

	_, err := c.Client.PutRolePolicy(context.TODO(), &iam.PutRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to attach policy to role")
		return fmt.Errorf("failed to attach policy to role '%s': %w", roleName, err)
	}

	logger.Info().Str("policy_name", policyName).Msg("attached policy to IAM role")
	return nil
}

func (c *IAMClient) createS3PolicyDocument(bucketNames []string) string {
	logger := c.logger.With().Strs("bucket_names", bucketNames).Logger()
	logger.Debug().Msg("creating S3 policy document")

	resources := make([]string, 0, len(bucketNames)*2)
	for _, bucket := range bucketNames {
		resources = append(resources,
			fmt.Sprintf("arn:aws:s3:::%s", bucket),
			fmt.Sprintf("arn:aws:s3:::%s/*", bucket))
	}

	policy := RolePolicyDocument{
		Version: "2012-10-17",
		Statement: []RoleStatement{
			{
				Effect: "Allow",
				Action: []string{
					"s3:ListBucket",
					"s3:GetObject",
					"s3:PutObject",
					"s3:GetBucketLocation",
				},
				Resource: resources,
			},
		},
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal policy")
		return ""
	}

	logger.Debug().RawJSON("policy", policyJSON).Msg("S3 policy document created")
	return string(policyJSON)
}
