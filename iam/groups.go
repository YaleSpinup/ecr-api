package iam

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
)

// GetGroup gets the details of an IAM group
func (i *IAM) GetGroupWithPath(ctx context.Context, name, path string) (*iam.Group, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	log.Infof("getting group %s with path %s", name, path)

	out, err := i.Service.GetGroupWithContext(ctx, &iam.GetGroupInput{
		GroupName: aws.String(name),
	})

	if err != nil {
		return nil, ErrCode("failed to get details for group", err)
	}

	log.Debugf("got output from get group %s: %+v", name, out)

	if aws.StringValue(out.Group.Path) != path {
		msg := fmt.Sprintf("group %s found, but not in path %s (actual path %s)", name, path, aws.StringValue(out.Group.Path))
		return nil, apierror.New(apierror.ErrNotFound, msg, nil)
	}

	return out.Group, nil
}

// CreateGroup handles creating an IAM group
func (i *IAM) CreateGroup(ctx context.Context, name, path string) (*iam.Group, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	log.Infof("creating group %s in path %s", name, path)

	out, err := i.Service.CreateGroupWithContext(ctx, &iam.CreateGroupInput{
		GroupName: aws.String(name),
		Path:      aws.String(path),
	})
	if err != nil {
		return nil, ErrCode("failed to create iam group", err)
	}

	log.Debugf("got output creating group %+v", out)

	return out.Group, nil
}

func (i *IAM) AttachGroupPolicy(ctx context.Context, groupName, policyArn string) error {
	if groupName == "" || policyArn == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("attaching policy %s to group %s", policyArn, groupName)

	_, err := i.Service.AttachGroupPolicyWithContext(ctx, &iam.AttachGroupPolicyInput{
		GroupName: aws.String(groupName),
		PolicyArn: aws.String(policyArn),
	})

	if err != nil {
		return ErrCode("failed to attach policy to group", err)
	}

	return nil
}

func (i *IAM) ListAttachedGroupPolicies(ctx context.Context, groupName, path string) ([]string, error) {
	if groupName == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	log.Infof("listing policies attached to group %s in path %s", groupName, path)

	policies := []string{}
	err := i.Service.ListAttachedGroupPoliciesPagesWithContext(ctx, &iam.ListAttachedGroupPoliciesInput{
		GroupName:  aws.String(groupName),
		PathPrefix: aws.String(path),
	}, func(page *iam.ListAttachedGroupPoliciesOutput, last bool) bool {
		for _, p := range page.AttachedPolicies {
			policies = append(policies, aws.StringValue(p.PolicyArn))
		}
		return true
	})

	if err != nil {
		return nil, ErrCode("failed to list attached group policies", err)
	}

	return policies, nil
}

// AddUserToGroup adds an existing user to an existing group
func (i *IAM) AddUserToGroup(ctx context.Context, userName, groupName string) error {
	if userName == "" || groupName == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("adding user %s to group %s", userName, groupName)

	if _, err := i.Service.AddUserToGroupWithContext(ctx, &iam.AddUserToGroupInput{
		UserName:  aws.String(userName),
		GroupName: aws.String(groupName),
	}); err != nil {
		return ErrCode("failed to add user to group", err)
	}

	return nil
}

// RemoveUserFromGroup removes a user from a group
func (i *IAM) RemoveUserFromGroup(ctx context.Context, userName, groupName string) error {
	if userName == "" || groupName == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("removing user %s from group %s", userName, groupName)

	if _, err := i.Service.RemoveUserFromGroupWithContext(ctx, &iam.RemoveUserFromGroupInput{
		UserName:  aws.String(userName),
		GroupName: aws.String(groupName),
	}); err != nil {
		return ErrCode("failed to remove user from group", err)
	}

	return nil
}
