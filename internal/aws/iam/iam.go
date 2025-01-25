// internal/aws/iam/iam.go
package iam

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

// IAMClient wraps the AWS IAM client
type IAMClient struct {
	Client *iam.Client
}

// NewIAMClient initializes a new IAM client
func NewIAMClient(cfg aws.Config) *IAMClient {
	return &IAMClient{
		Client: iam.NewFromConfig(cfg),
	}
}

// SetupTrustRelationship creates or updates an IAM role with a trust relationship and attaches necessary policies
func (c *IAMClient) SetupTrustRelationship(roleName, trustedAccountID, externalID string, workspace string,
	workergroup string, action string, bucketNames []string) error {
	if err := c.validateInputs(roleName, trustedAccountID, workspace, workergroup, bucketNames); err != nil {
		return err
	}

	trustPolicy := c.createTrustPolicy(trustedAccountID, workspace, workergroup, externalID, action)

	if err := c.ensureRoleExists(roleName, trustPolicy); err != nil {
		return err
	}

	if err := c.attachS3Policies(roleName, bucketNames); err != nil {
		return err
	}

	return nil
}

func (c *IAMClient) validateInputs(roleName, trustedAccountID, workspace, workergroup string, bucketNames []string) error {
	if roleName == "" {
		return fmt.Errorf("roleName cannot be empty")
	}
	if trustedAccountID == "" {
		return fmt.Errorf("trustedAccountID cannot be empty")
	}
	if len(bucketNames) == 0 {
		return fmt.Errorf("at least one bucketName must be provided")
	}
	return nil
}

func (c *IAMClient) createTrustPolicy(trustedAccountID, workspace, workergroup, externalID, action string) string {
	if action != "search" {
		return fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [{
                "Effect": "Allow",
                "Principal": { "AWS": "arn:aws:iam::%s:role/%s-%s" },
                "Action": ["sts:AssumeRole","sts:TagSession","sts:SetSourceIdentity"],
                "Condition": {
                    "StringEquals": { "sts:ExternalId": "%s" }
                }
            }]
        }`, trustedAccountID, workspace, workergroup, externalID)
	}

	return fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Principal": { "AWS": "arn:aws:iam::%s:role/search-exec-main" },
            "Action": ["sts:AssumeRole","sts:TagSession","sts:SetSourceIdentity"],
            "Condition": {
                "StringEquals": { "sts:ExternalId": "%s" }
            }
        }]
    }`, trustedAccountID, externalID)
}

func (c *IAMClient) ensureRoleExists(roleName, trustPolicy string) error {
	_, err := c.Client.GetRole(context.TODO(), &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})

	if err != nil {
		var noSuchEntity *types.NoSuchEntityException
		if errors.As(err, &noSuchEntity) {
			return c.createRole(roleName, trustPolicy)
		}
		return fmt.Errorf("failed to get IAM role: %w", err)
	}

	return c.updateRoleTrustPolicy(roleName, trustPolicy)
}

func (c *IAMClient) createRole(roleName, trustPolicy string) error {
	_, err := c.Client.CreateRole(context.TODO(), &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(trustPolicy),
		Description:              aws.String("Role for cross-account access to S3"),
	})
	if err != nil {
		return fmt.Errorf("failed to create IAM role: %w", err)
	}
	fmt.Printf("Created IAM role '%s' with trust relationship.\n", roleName)
	return nil
}

func (c *IAMClient) updateRoleTrustPolicy(roleName, trustPolicy string) error {
	_, err := c.Client.UpdateAssumeRolePolicy(context.TODO(), &iam.UpdateAssumeRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyDocument: aws.String(trustPolicy),
	})
	if err != nil {
		return fmt.Errorf("failed to update trust policy for role '%s': %w", roleName, err)
	}
	fmt.Printf("Updated trust relationship for IAM role '%s'.\n", roleName)
	return nil
}

func (c *IAMClient) attachS3Policies(roleName string, bucketNames []string) error {
	policyName := "CrossAccountAccessPolicy"
	policyDocument := c.createS3PolicyDocument(bucketNames)

	_, err := c.Client.PutRolePolicy(context.TODO(), &iam.PutRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
	})
	if err != nil {
		return fmt.Errorf("failed to attach policy to role '%s': %w", roleName, err)
	}
	fmt.Printf("Attached policy '%s' to IAM role '%s'.\n", policyName, roleName)
	return nil
}

func (c *IAMClient) createS3PolicyDocument(bucketNames []string) string {
	var bucketArns []string
	for _, bucket := range bucketNames {
		bucketArns = append(bucketArns, fmt.Sprintf("\"arn:aws:s3:::%s\"", bucket))
		bucketArns = append(bucketArns, fmt.Sprintf("\"arn:aws:s3:::%s/*\"", bucket))
	}

	return fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetObject",
                "s3:PutObject",
                "s3:GetBucketLocation"
            ],
            "Resource": [%s]
        }]
    }`, strings.Join(bucketArns, ", "))
}
