package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"

	im "github.com/YaleSpinup/ecr-api/iam"
	log "github.com/sirupsen/logrus"
)

var ecrAdminPolicyDoc string
var EcrAdminPolicy = im.PolicyDocument{
	Version: "2012-10-17",
	Statement: []im.StatementEntry{
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
			Resource: "*",
			Condition: im.Condition{
				"StringEquals": im.ConditionStatement{
					"aws:ResourceTag/spinup:org":     "${aws:PrincipalTag/spinup:org}",
					"aws:ResourceTag/spinup:spaceid": "${aws:PrincipalTag/spinup:spaceid}",
					"aws:ResourceTag/Name":           "${aws:PrincipalTag/Name}",
				},
			},
		},
		{
			Sid:      "AllowDockerLogin",
			Effect:   "Allow",
			Action:   []string{"ecr:GetAuthorizationToken"},
			Resource: "*",
		},
	},
}

type iamOrchestrator struct {
	client im.IAM
	org    string
}

func (o *iamOrchestrator) listRepositoryUsers(ctx context.Context, group, name string) ([]string, error) {

	path := fmt.Sprintf("/spinup/%s", o.org)

	if group != "" {
		path = path + fmt.Sprintf("/%s", group)
	}

	if name != "" {
		path = path + fmt.Sprintf("/%s", name)
	}

	users, err := o.client.ListUsers(ctx, path)
	if err != nil {
		return nil, err
	}

	ps := strings.Split(path, "/")
	if len(ps) > 2 {
		prefix := fmt.Sprintf("%s-%s-", ps[len(ps)-2], ps[len(ps)-1])

		trimmed := make([]string, 0, len(users))
		for _, u := range users {
			log.Debugf("trimming prefix '%s' from username %s", prefix, u)
			u = strings.TrimPrefix(u, prefix)
			trimmed = append(trimmed, u)
		}
		users = trimmed
	}

	return users, nil
}

func (o *iamOrchestrator) getRepositoryUser(ctx context.Context, group, name, user string) (*RepositoryUserResponse, error) {

	path := fmt.Sprintf("/spinup/%s/%s/%s/", o.org, group, name)
	userName := fmt.Sprintf("%s-%s-%s", group, name, user)

	iamUser, err := o.client.GetUserWithPath(ctx, path, userName)
	if err != nil {
		return nil, err
	}

	keys, err := o.client.ListAccessKeys(ctx, userName)
	if err != nil {
		return nil, err
	}

	return repositoryUserResponseFromIAM(iamUser, keys), nil
}

func (o *iamOrchestrator) repositoryUserDelete(ctx context.Context, name, group, user string) error {
	path := fmt.Sprintf("/spinup/%s/%s/%s/", o.org, group, name)
	userName := fmt.Sprintf("%s-%s-%s", group, name, user)

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

	if err := o.client.DeleteUser(ctx, userName); err != nil {
		return err
	}

	return nil
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

	doc := im.PolicyDocument{}
	if err := json.Unmarshal([]byte(d), &doc); err != nil {
		return "", err
	}

	if !awsutil.DeepEqual(doc, EcrAdminPolicy) {
		log.Warn("policy document is not the same, updating")

		if err := o.client.UpdatePolicy(ctx, aws.StringValue(policy.Arn), ecrAdminPolicyDoc); err != nil {
			return "", err
		}

		// TODO: delete old version, only 5 versions allowed
	}

	return aws.StringValue(policy.Arn), nil
}

func (o *iamOrchestrator) userCreateGroupIfMissing(ctx context.Context, name, path, policyArn string) error {
	log.Infof("creating group %s in %s and assigning policy %s if missing", name, path, policyArn)

	_, err := o.client.GetGroupWithPath(ctx, name, path)
	if err != nil {
		if aerr, ok := err.(apierror.Error); ok && aerr.Code == apierror.ErrNotFound {
			log.Infof("policy %s not found, creating", name)

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

func (o *iamOrchestrator) repositoryUserCreate(ctx context.Context, name, group, groupName string, req RepositoryUserCreateRequest) (*RepositoryUserResponse, error) {
	log.Infof("creating repository %s user %s in group %s in iam group %s", name, req.UserName, group, groupName)

	path := fmt.Sprintf("/spinup/%s/%s/%s/", o.org, group, name)
	userName := fmt.Sprintf("%s-%s-%s", group, name, req.UserName)
	repository := fmt.Sprintf("%s/%s", group, name)

	req.Tags = normalizeTags(o.org, group, repository, req.Tags)

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

	return repositoryUserResponseFromIAM(user, nil), nil
}