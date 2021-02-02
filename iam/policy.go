package iam

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
)

func (i *IAM) GetPolicyByName(ctx context.Context, name, path string) (*iam.Policy, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	var policy *iam.Policy
	err := i.Service.ListPoliciesPagesWithContext(ctx, &iam.ListPoliciesInput{PathPrefix: aws.String(path)}, func(page *iam.ListPoliciesOutput, lastPage bool) bool {
		for _, p := range page.Policies {
			if aws.StringValue(p.PolicyName) == name {
				policy = p
				log.Debugf("found policy with name %s and path %s", name, path)
				return false
			}
		}
		return true
	})

	if err != nil {
		return nil, ErrCode("failed finding policy", err)
	}

	if policy == nil {
		return nil, apierror.New(apierror.ErrNotFound, "policy not found", nil)
	}

	if aws.StringValue(policy.Path) != path {
		msg := fmt.Sprintf("policy %s found, but not in path %s (actual path %s)", name, path, aws.StringValue(policy.Path))
		return nil, apierror.New(apierror.ErrNotFound, msg, nil)
	}

	return policy, nil
}

func (i *IAM) GetDefaultPolicyVersion(ctx context.Context, arn, version string) (*iam.PolicyVersion, error) {
	if arn == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting policy document for %s", arn)

	out, err := i.Service.GetPolicyVersionWithContext(ctx, &iam.GetPolicyVersionInput{
		PolicyArn: aws.String(arn),
		VersionId: aws.String(version),
	})

	if err != nil {
		return nil, ErrCode("failed to get policy version", err)
	}

	log.Debugf("got output from getting policy version: %+v", out)

	return out.PolicyVersion, nil
}

func (i *IAM) WaitForPolicy(ctx context.Context, policyArn string) error {
	if policyArn == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if err := i.Service.WaitUntilPolicyExistsWithContext(ctx, &iam.GetPolicyInput{
		PolicyArn: aws.String(policyArn),
	}); err != nil {
		return ErrCode("failed waiting for policy to create", err)
	}

	return nil
}

func (i *IAM) CreatePolicy(ctx context.Context, name, path, policyDoc string) (*iam.Policy, error) {
	if name == "" || policyDoc == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	log.Infof("creating policy %s in path %s", name, path)

	out, err := i.Service.CreatePolicyWithContext(ctx, &iam.CreatePolicyInput{
		Path:           aws.String(path),
		PolicyDocument: aws.String(policyDoc),
		PolicyName:     aws.String(name),
	})

	if err != nil {
		return nil, ErrCode("failed to create policy", err)
	}

	log.Debugf("got output from create policy:  %+v", out)

	return out.Policy, nil
}

func (i *IAM) UpdatePolicy(ctx context.Context, arn, policyDoc string) error {
	if arn == "" || policyDoc == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("updating policy %s by creating new version and setting as default", arn)

	out, err := i.Service.CreatePolicyVersionWithContext(ctx, &iam.CreatePolicyVersionInput{
		PolicyArn:      aws.String(arn),
		PolicyDocument: aws.String(policyDoc),
		SetAsDefault:   aws.Bool(true),
	})

	if err != nil {
		return ErrCode("failed to update policy", err)
	}

	log.Debugf("got output from create policy version: %+v", out)

	return nil
}
