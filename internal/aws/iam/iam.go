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
    workergroup string, action string,  bucketNames []string) error {
    // Validate input parameters
    if roleName == "" {
        return fmt.Errorf("roleName cannot be empty")
    }
    if trustedAccountID == "" {
        return fmt.Errorf("trustedAccountID cannot be empty")
    }
    if len(bucketNames) == 0 {
        return fmt.Errorf("at least one bucketName must be provided")
    }
    if workergroup == "" && action != "search" {
        workergroup := "default"
    }
    if workspace == "" && action != "search" {
        workspace := "main"
    }

    if action != "search" {
    // Define trust relationship policy
    trustPolicy := fmt.Sprintf(`{
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
    } else {
        trustPolicy := fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [{
                "Effect": "Allow",
                "Principal": { "AWS": "arn:aws:iam::%s:role/search-exec-*" },
                "Action": ["sts:AssumeRole","sts:TagSession","sts:SetSourceIdentity"],
                "Condition": {
                    "StringEquals": { "sts:ExternalId": "%s" }
                }
            }]
        }`, trustedAccountID, externalID)
    }


    // Check if the role exists
    _, err := c.Client.GetRole(context.TODO(), &iam.GetRoleInput{
        RoleName: aws.String(roleName),
    })

    if err != nil {
        var noSuchEntity *types.NoSuchEntityException
        if errors.As(err, &noSuchEntity) {
            // Create the role
            _, err = c.Client.CreateRole(context.TODO(), &iam.CreateRoleInput{
                RoleName:                 aws.String(roleName),
                AssumeRolePolicyDocument: aws.String(trustPolicy),
                Description:              aws.String("Role for cross-account access to S3"),
            })
            if err != nil {
                return fmt.Errorf("failed to create IAM role: %w", err)
            }
            fmt.Printf("Created IAM role '%s' with trust relationship.\n", roleName)
        } else {
            return fmt.Errorf("failed to get IAM role: %w", err)
        }
    } else {
        // Update the trust policy
        _, err = c.Client.UpdateAssumeRolePolicy(context.TODO(), &iam.UpdateAssumeRolePolicyInput{
            RoleName:       aws.String(roleName),
            PolicyDocument: aws.String(trustPolicy),
        })
        if err != nil {
            return fmt.Errorf("failed to update trust policy for role '%s': %w", roleName, err)
        }
        fmt.Printf("Updated trust relationship for IAM role '%s'.\n", roleName)
    }

    // Attach policies as needed (e.g., permissions for S3 access)
    // For example, attach an inline policy granting S3 access

    policyName := "CrossAccountAccessPolicy"

    // Define the policy document with the provided bucket names
    // Generate ARN strings for each bucket
    var bucketArns []string
    for _, bucket := range bucketNames {
        bucketArns = append(bucketArns, fmt.Sprintf("\"arn:aws:s3:::%s\"", bucket))
        bucketArns = append(bucketArns, fmt.Sprintf("\"arn:aws:s3:::%s/*\"", bucket))
    }
    policyDocument := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetObject",
                "s3:PutObject"
            ],
            "Resource": [%s]
        }]
    }`, strings.Join(bucketArns, ", "))

    // Attach the policy to the role
    _, err = c.Client.PutRolePolicy(context.TODO(), &iam.PutRolePolicyInput{
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