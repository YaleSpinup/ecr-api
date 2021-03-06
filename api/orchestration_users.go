package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/YaleSpinup/ecr-api/iam"
	log "github.com/sirupsen/logrus"
)

var ecrAdminPolicyDoc string
var EcrAdminPolicy = iam.PolicyDocument{
	Version: "2012-10-17",
	Statement: []iam.StatementEntry{
		{
			Sid:    "AllowActionsOnRepositoriesInSpaceAndOrg",
			Effect: "Allow",
			Action: []string{
				"ecr:PutLifecyclePolicy",
				"ecr:PutImageTagMutability",
				"ecr:DescribeImageScanFindings",
				"ecr:GetDownloadUrlForLayer",
				"ecr:GetAuthorizationToken",
				"ecr:UploadLayerPart",
				"ecr:BatchDeleteImage",
				"ecr:ListImages",
				"ecr:DeleteLifecyclePolicy",
				"ecr:PutImage",
				"ecr:BatchGetImage",
				"ecr:CompleteLayerUpload",
				"ecr:DescribeImages",
				"ecr:DeleteRegistryPolicy",
				"ecr:InitiateLayerUpload",
				"ecr:BatchCheckLayerAvailability",
			},
			Resource: []string{"*"},
			Condition: iam.Condition{
				"StringEqualsIgnoreCase": iam.ConditionStatement{
					"aws:ResourceTag/Name":           []string{"${aws:PrincipalTag/ResourceName}"},
					"aws:ResourceTag/spinup:org":     []string{"${aws:PrincipalTag/spinup:org}"},
					"aws:ResourceTag/spinup:spaceid": []string{"${aws:PrincipalTag/spinup:spaceid}"},
				},
			},
		},
		{
			Sid:      "AllowDockerLogin",
			Effect:   "Allow",
			Action:   []string{"ecr:GetAuthorizationToken"},
			Resource: []string{"*"},
		},
	},
}

// listRepositoryUsers lists users in a repository
func (o *iamOrchestrator) listRepositoryUsers(ctx context.Context, group, name string) ([]string, error) {
	path := fmt.Sprintf("/spinup/%s/%s/%s", o.org, group, name)

	users, err := o.client.ListUsers(ctx, path)
	if err != nil {
		return nil, err
	}

	prefix := name + "-"

	trimmed := make([]string, 0, len(users))
	for _, u := range users {
		log.Debugf("trimming prefix '%s' from username %s", prefix, u)
		u = strings.TrimPrefix(u, prefix)
		trimmed = append(trimmed, u)
	}
	users = trimmed

	return users, nil
}

// getRepositoryUser gets the details about a user
func (o *iamOrchestrator) getRepositoryUser(ctx context.Context, group, name, user string) (*RepositoryUserResponse, error) {
	path := fmt.Sprintf("/spinup/%s/%s/%s/", o.org, group, name)
	userName := fmt.Sprintf("%s-%s", name, user)

	iamUser, err := o.client.GetUserWithPath(ctx, path, userName)
	if err != nil {
		return nil, err
	}

	keys, err := o.client.ListAccessKeys(ctx, userName)
	if err != nil {
		return nil, err
	}

	groups, err := o.client.ListGroupsForUser(ctx, userName)
	if err != nil {
		return nil, err
	}

	return repositoryUserResponseFromIAM(o.org, iamUser, keys, groups), nil
}

// repositoryUserDelete orchestrates removing a user from all groups and deleting the user
func (o *iamOrchestrator) repositoryUserDelete(ctx context.Context, name, group, user string) error {
	path := fmt.Sprintf("/spinup/%s/%s/%s/", o.org, group, name)
	userName := fmt.Sprintf("%s-%s", name, user)

	if _, err := o.client.GetUserWithPath(ctx, path, userName); err != nil {
		return err
	}

	groups, err := o.client.ListGroupsForUser(ctx, userName)
	if err != nil {
		return err
	}

	for _, g := range groups {
		if err := o.client.RemoveUserFromGroup(ctx, userName, g); err != nil {
			return err
		}
	}

	keys, err := o.client.ListAccessKeys(ctx, userName)
	if err != nil {
		return err
	}

	for _, k := range keys {
		if err := o.client.DeleteAccessKey(ctx, userName, aws.StringValue(k.AccessKeyId)); err != nil {
			return err
		}
	}

	if err := o.client.DeleteUser(ctx, userName); err != nil {
		return err
	}

	return nil
}

// repositoryUserDeleteAll deletes all users for a repository
func (o *iamOrchestrator) repositoryUserDeleteAll(ctx context.Context, name, group string) ([]string, error) {

	// list all users for the repository
	users, err := o.listRepositoryUsers(ctx, group, name)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if err := o.repositoryUserDelete(ctx, name, group, u); err != nil {
			log.Errorf("failed to delete repository %s user %s: %s", name, u, err)
		}
	}

	return users, nil
}

// prepareAccount sets up the account for user management by creating the ECR admin policy and group
func (o *iamOrchestrator) prepareAccount(ctx context.Context) (string, error) {
	log.Info("preparing account for user management")

	path := fmt.Sprintf("/spinup/%s/", o.org)

	policyName := fmt.Sprintf("SpinupECRAdminPolicy-%s", o.org)
	policyArn, err := o.userCreatePolicyIfMissing(ctx, policyName, path)
	if err != nil {
		return "", err
	}

	groupName := fmt.Sprintf("SpinupECRAdminGroup-%s", o.org)
	if err := o.userCreateGroupIfMissing(ctx, groupName, path, policyArn); err != nil {
		return "", err
	}

	return groupName, err
}

// userCreatePolicyIfMissing gets the given policy by name.  if the policy isn't found it simply creates the policy and
// returns.  if the policy is found, it gets the policy document and compares to the expected policy document, updating
// if they differ.
func (o *iamOrchestrator) userCreatePolicyIfMissing(ctx context.Context, name, path string) (string, error) {
	log.Infof("creating policy %s in %s if missing", name, path)

	policy, err := o.client.GetPolicyByName(ctx, name, path)
	if err != nil {
		if aerr, ok := err.(apierror.Error); ok && aerr.Code == apierror.ErrNotFound {
			log.Infof("policy %s not found, creating", name)
		} else {
			return "", err
		}
	}

	// if the policy isn't found, create it and return
	if policy == nil {
		out, err := o.client.CreatePolicy(ctx, name, path, ecrAdminPolicyDoc)
		if err != nil {
			return "", err
		}

		if err := o.client.WaitForPolicy(ctx, aws.StringValue(out.Arn)); err != nil {
			return "", err
		}

		return aws.StringValue(out.Arn), nil
	}

	out, err := o.client.GetDefaultPolicyVersion(ctx, aws.StringValue(policy.Arn), aws.StringValue(policy.DefaultVersionId))
	if err != nil {
		return "", err
	}

	// Document is returned url encoded, we must decode it to unmarshal and compare
	d, err := url.QueryUnescape(aws.StringValue(out.Document))
	if err != nil {
		return "", err
	}

	// If we cannot unmarshal the document we received into an iam.PolicyDocument or if
	// the document doesn't match, let's try to update it.  If unmarshaling fails, we assume
	// our struct has changed (for example going from Resource string to Resource []string)
	var updatePolicy bool
	doc := iam.PolicyDocument{}
	if err := json.Unmarshal([]byte(d), &doc); err != nil {
		log.Warnf("error getting policy document: %s, updating", err)
		updatePolicy = true
	} else if !iam.PolicyDeepEqual(doc, EcrAdminPolicy) {
		log.Warn("policy document is not the same, updating")
		updatePolicy = true
	}

	if updatePolicy {
		if err := o.client.UpdatePolicy(ctx, aws.StringValue(policy.Arn), ecrAdminPolicyDoc); err != nil {
			return "", err
		}

		// TODO: delete old version, only 5 versions allowed
	}

	return aws.StringValue(policy.Arn), nil
}

func (o *iamOrchestrator) userCreateGroupIfMissing(ctx context.Context, name, path, policyArn string) error {
	log.Infof("creating group %s in %s and assigning policy %s if missing", name, path, policyArn)

	if _, err := o.client.GetGroupWithPath(ctx, name, path); err != nil {
		if aerr, ok := err.(apierror.Error); ok && aerr.Code == apierror.ErrNotFound {
			log.Infof("group %s not found, creating", name)

			if _, err := o.client.CreateGroup(ctx, name, path); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	attachedPolicies, err := o.client.ListAttachedGroupPolicies(ctx, name, path)
	if err != nil {
		return err
	}

	// return if the policy is already attached to the group
	for _, p := range attachedPolicies {
		if p != policyArn {
			continue
		}

		return nil
	}

	if err := o.client.AttachGroupPolicy(ctx, name, policyArn); err != nil {
		return err
	}

	return nil
}

func (o *iamOrchestrator) repositoryUserCreate(ctx context.Context, name, group, groupName string, req *RepositoryUserCreateRequest) (*RepositoryUserResponse, error) {
	log.Infof("creating repository %s user %s in group %s in iam group %s", name, req.UserName, group, groupName)

	path := fmt.Sprintf("/spinup/%s/%s/%s/", o.org, group, name)
	userName := fmt.Sprintf("%s-%s", name, req.UserName)
	repository := fmt.Sprintf("%s/%s", group, name)

	req.Tags = normalizeUserTags(o.org, group, repository, userName, req.Tags)

	user, err := o.client.CreateUser(ctx, userName, path, toIAMTags(req.Tags))
	if err != nil {
		return nil, err
	}

	if err := o.client.WaitForUser(ctx, userName); err != nil {
		return nil, err
	}

	// append the org to the passed group(s) and add user to the group
	// TODO rollback on failure
	for _, g := range req.Groups {
		grp := fmt.Sprintf("%s-%s", g, o.org)

		if err := o.client.AddUserToGroup(ctx, userName, grp); err != nil {
			return nil, err
		}
	}

	return repositoryUserResponseFromIAM(o.org, user, nil, []string{groupName}), nil
}

func (o *iamOrchestrator) repositoryUserUpdate(ctx context.Context, name, group, uname string, req *RepositoryUserUpdateRequest) (*RepositoryUserResponse, error) {
	log.Infof("updating repository %s user %s in group %s", name, uname, group)

	userName := fmt.Sprintf("%s-%s", name, uname)
	repository := fmt.Sprintf("%s/%s", group, name)

	response := &RepositoryUserResponse{
		UserName: uname,
	}

	if req.Tags != nil {
		req.Tags = normalizeUserTags(o.org, group, repository, userName, req.Tags)
		if err := o.client.TagUser(ctx, userName, toIAMTags(req.Tags)); err != nil {
			return nil, err
		}
		response.Tags = req.Tags
	}

	if req.ResetKey {
		// get a list of users access keys
		keys, err := o.client.ListAccessKeys(ctx, userName)
		if err != nil {
			return nil, err
		}

		newKeyOut, err := o.client.CreateAccessKey(ctx, userName)
		if err != nil {
			return nil, err
		}
		response.AccessKey = newKeyOut

		deletedKeyIds := make([]string, 0, len(keys))
		// delete the old access keys
		for _, k := range keys {
			err = o.client.DeleteAccessKey(ctx, userName, aws.StringValue(k.AccessKeyId))
			if err != nil {
				return response, err
			}
			deletedKeyIds = append(deletedKeyIds, aws.StringValue(k.AccessKeyId))
		}

		response.DeletedAccessKeys = deletedKeyIds
	}

	return response, nil
}
