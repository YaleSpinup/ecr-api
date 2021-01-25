package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/ecr"

	"github.com/YaleSpinup/apierror"
	im "github.com/YaleSpinup/ecr-api/iam"
	log "github.com/sirupsen/logrus"
)

// repositoryCreate orchestrates the creation of a repository from the RepositoryCreateRequest
func (e *ecrOrchestrator) repositoryCreate(ctx context.Context, account, group string, req *RepositoryCreateRequest) (*RepositoryResponse, error) {
	repository := fmt.Sprintf("%s/%s", group, req.RepositoryName)

	log.Debugf("creating %s repository with request %+v", repository, req)

	req.Tags = normalizeTags(e.org, group, repository, req.Tags)

	scanOnPush := false
	if req.ScanOnPush != "" {
		b, err := strconv.ParseBool(req.ScanOnPush)
		if err != nil {
			return nil, err
		}
		scanOnPush = b
	}

	input := &ecr.CreateRepositoryInput{
		EncryptionConfiguration: &ecr.EncryptionConfiguration{
			EncryptionType: aws.String("AES256"),
		},
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(scanOnPush),
		},
		RepositoryName: aws.String(repository),
		Tags:           toECRTags(req.Tags),
	}

	if req.KmsKeyId != "" {
		input = input.SetEncryptionConfiguration(&ecr.EncryptionConfiguration{
			EncryptionType: aws.String("KMS"),
			KmsKey:         aws.String(req.KmsKeyId),
		})
	}

	log.Debugf("creating repository with input %s", awsutil.Prettify(input))

	out, err := e.client.CreateRepository(ctx, input)
	if err != nil {
		return nil, err
	}

	tags, err := e.client.GetRepositoryTags(ctx, aws.StringValue(out.RepositoryArn))
	if err != nil {
		return nil, err
	}

	log.Debugf("got output %+v", out)

	return repositoryResponseFromECR(out, tags), nil
}

// repositoryDelete orchestrates the deletion of a repository
func (e *ecrOrchestrator) repositoryDelete(ctx context.Context, account, group, name string) (*RepositoryResponse, error) {
	repository := fmt.Sprintf("%s/%s", group, name)

	log.Debugf("deleting repository %s", repository)

	repo, err := e.client.GetRepositories(ctx, repository)
	if err != nil {
		return nil, err
	}

	tags, err := e.client.GetRepositoryTags(ctx, aws.StringValue(repo.RepositoryArn))
	if err != nil {
		return nil, err
	}

	out, err := e.client.DeleteRepository(ctx, repository)
	if err != nil {
		return nil, err
	}

	log.Debugf("got output %+v", out)

	return repositoryResponseFromECR(out, tags), nil
}

// repositoryUpdate orchestrates updating a repository
func (e *ecrOrchestrator) repositoryUpdate(ctx context.Context, account, group, name string, req *RepositoryUpdateRequest) (*RepositoryResponse, error) {
	repository := fmt.Sprintf("%s/%s", group, name)

	log.Debugf("updating %s repository with request %+v", repository, req)

	req.Tags = normalizeTags(e.org, group, repository, req.Tags)

	repo, err := e.client.GetRepositories(ctx, repository)
	if err != nil {
		return nil, err
	}

	if req.ScanOnPush != "" {
		scanOnPush, err := strconv.ParseBool(req.ScanOnPush)
		if err != nil {
			return nil, err
		}

		if err := e.client.SetImageScanningConfiguration(ctx, repository, scanOnPush); err != nil {
			return nil, err
		}
	}

	if req.Tags != nil {
		if err := e.client.UpdateRepositoryTags(ctx, aws.StringValue(repo.RepositoryArn), toECRTags(req.Tags)); err != nil {
			return nil, err
		}
	}

	repo, err = e.client.GetRepositories(ctx, repository)
	if err != nil {
		return nil, err
	}

	tags, err := e.client.GetRepositoryTags(ctx, aws.StringValue(repo.RepositoryArn))
	if err != nil {
		return nil, err
	}

	return repositoryResponseFromECR(repo, tags), nil
}

// userCreatePolicyIfMissing gets the given policy by name.  if the policy isn't found it simply creates the policy and
// returns.  if the policy is found, it gets the policy document and compares to the expected policy document, updating
// if they differ.
func (i *iamOrchestrator) userCreatePolicyIfMissing(ctx context.Context, name, path string) (string, error) {
	log.Debugf("creating policy %s in %s if missing", name, path)

	policy, err := i.client.GetPolicyByName(ctx, name, path)
	if err != nil {
		if aerr, ok := err.(apierror.Error); ok && aerr.Code == apierror.ErrNotFound {
			log.Infof("policy %s not found, creating", name)
		} else {
			return "", err
		}
	}

	// if the policy isn't found, create it and return
	if policy == nil {
		out, err := i.client.CreatePolicy(ctx, name, path, ecrAdminPolicyDoc)
		if err != nil {
			return "", err
		}

		if err := i.client.WaitForPolicy(ctx, aws.StringValue(out.Arn)); err != nil {
			return "", err
		}

		return aws.StringValue(out.Arn), nil
	}

	out, err := i.client.GetDefaultPolicyVersion(ctx, aws.StringValue(policy.Arn), aws.StringValue(policy.DefaultVersionId))
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

		if err := i.client.UpdatePolicy(ctx, aws.StringValue(policy.Arn), ecrAdminPolicyDoc); err != nil {
			return "", err
		}

		// TODO: delete old version, only 5 versions allowed
	}

	return aws.StringValue(policy.Arn), nil
}

func (i *iamOrchestrator) userCreateGroupIfMissing(ctx context.Context, name, path, policyArn string) error {
	log.Debugf("creating group %s in %s and assigning policy %s if missing", name, path, policyArn)

	_, err := i.client.GetGroupWithPath(ctx, name, path)
	if err != nil {
		if aerr, ok := err.(apierror.Error); ok && aerr.Code == apierror.ErrNotFound {
			log.Infof("policy %s not found, creating", name)

			if _, err := i.client.CreateGroup(ctx, name, path); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	attachedPolicies, err := i.client.ListAttachedGroupPolicies(ctx, name, path)
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

	if err := i.client.AttachGroupPolicy(ctx, name, policyArn); err != nil {
		return err
	}

	return nil
}

func (i *iamOrchestrator) repositoryUserCreate(ctx context.Context, name, group, groupName string, req RepositoryUserCreateRequest) (*RepositoryUserResponse, error) {
	path := fmt.Sprintf("/spinup/%s/%s/%s/", i.org, group, name)
	userName := fmt.Sprintf("%s-%s-%s", group, name, req.UserName)
	repository := fmt.Sprintf("%s/%s", group, name)

	req.Tags = normalizeTags(i.org, group, repository, req.Tags)

	user, err := i.client.CreateUser(ctx, userName, path, toIAMTags(req.Tags))
	if err != nil {
		return nil, err
	}

	if err := i.client.WaitForUser(ctx, userName); err != nil {
		return nil, err
	}

	// append the org to the passed group(s) and add user to the group
	// TODO rollback on failure
	for _, g := range req.Groups {
		grp := fmt.Sprintf("%s-%s", g, i.org)

		if err := i.client.AddUserToGroup(ctx, userName, grp); err != nil {
			return nil, err
		}
	}

	return repositoryUserResponseFromIAM(user, nil), nil
}

func (i *iamOrchestrator) repositoryUserDelete(ctx context.Context, name, group, user string) error {
	path := fmt.Sprintf("/spinup/%s/%s/%s/", i.org, group, name)
	userName := fmt.Sprintf("%s-%s-%s", group, name, user)

	if _, err := i.client.GetUserWithPath(ctx, path, userName); err != nil {
		return err
	}

	groups, err := i.client.ListGroupsForUser(ctx, userName)
	if err != nil {
		return err
	}

	for _, g := range groups {
		if err := i.client.RemoveUserFromGroup(ctx, userName, g); err != nil {
			return err
		}
	}

	if err := i.client.DeleteUser(ctx, userName); err != nil {
		return err
	}

	return nil
}
