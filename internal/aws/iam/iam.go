// internal/aws/iam/iam.go
package iam

import (
    "context"
    "errors"
    "fmt"

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
func (c *IAMClient) SetupTrustRelationship(roleName, trustedAccountID, externalID string) error {
    // Define trust relationship policy
    trustPolicy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Principal": { "AWS": "arn:aws:iam::%s:root" },
            "Action": "sts:AssumeRole",
            "Condition": {
                "StringEquals": { "sts:ExternalId": "%s" }
            }
        }]
    }`, trustedAccountID, externalID)

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

    // Define the policy document
    policyDocument := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetObject",
                "s3:PutObject"
            ],
            "Resource": [
                "arn:aws:s3:::%s",
                "arn:aws:s3:::%s/*"
            ]
        }]
    }`, "your-s3-bucket-name", "your-s3-bucket-name") // Replace with actual bucket names

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
