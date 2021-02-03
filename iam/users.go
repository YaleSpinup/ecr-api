package iam

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
)

// ListUsers lists all of the users in a path prefix, "/" by default
func (i *IAM) ListUsers(ctx context.Context, path string) ([]string, error) {
	if path == "" {
		path = "/"
	}

	out, err := i.Service.ListUsersWithContext(ctx, &iam.ListUsersInput{
		PathPrefix: aws.String(path),
	})

	if err != nil {
		return nil, ErrCode("failed to list users", err)
	}

	log.Debugf("got output from list users: %+v", out)

	users := make([]string, 0, len(out.Users))
	for _, u := range out.Users {
		users = append(users, aws.StringValue(u.UserName))
	}

	return users, nil
}

// GetUserWithPath gets details about a user and returns an error if the path doesn't match
func (i *IAM) GetUserWithPath(ctx context.Context, path, name string) (*iam.User, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	log.Infof("getting user %s with path %s", name, path)

	out, err := i.Service.GetUserWithContext(ctx, &iam.GetUserInput{
		UserName: aws.String(name),
	})

	if err != nil {
		return nil, ErrCode("failed to get user", err)
	}

	log.Debugf("got output from get user: %+v", out)

	if aws.StringValue(out.User.Path) != path {
		msg := fmt.Sprintf("user %s found, but not in path %s (actual path %s)", name, path, aws.StringValue(out.User.Path))
		return nil, apierror.New(apierror.ErrNotFound, msg, nil)
	}

	return out.User, nil
}

func (i *IAM) CreateAccessKey(ctx context.Context, name string) (*iam.AccessKey, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating access key for %s", name)

	out, err := i.Service.CreateAccessKeyWithContext(ctx, &iam.CreateAccessKeyInput{
		UserName: aws.String(name),
	})

	if err != nil {
		return nil, ErrCode("failed to create access keys", err)
	}

	log.Debugf("got output from create access keys: %+v", out)

	return out.AccessKey, nil
}

func (i *IAM) DeleteAccessKey(ctx context.Context, name, keyId string) error {
	if name == "" || keyId == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("deleting access key %s for %s", keyId, name)

	_, err := i.Service.DeleteAccessKeyWithContext(ctx, &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(keyId),
		UserName:    aws.String(name),
	})

	if err != nil {
		return ErrCode("failed to delete access keys", err)
	}

	return nil
}

func (i *IAM) ListAccessKeys(ctx context.Context, name string) ([]*iam.AccessKeyMetadata, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("listing access keys for %s", name)

	out, err := i.Service.ListAccessKeysWithContext(ctx, &iam.ListAccessKeysInput{
		UserName: aws.String(name),
	})

	if err != nil {
		return nil, ErrCode("failed to list access keys", err)
	}

	log.Debugf("got output from list access keys: %+v", out)

	return out.AccessKeyMetadata, nil
}

func (i *IAM) CreateUser(ctx context.Context, name, path string, tags []*iam.Tag) (*iam.User, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	log.Infof("creating user %s in path %s", name, path)

	out, err := i.Service.CreateUserWithContext(ctx, &iam.CreateUserInput{
		UserName: aws.String(name),
		Path:     aws.String(path),
		Tags:     tags,
	})

	if err != nil {
		return nil, ErrCode("failed to create user", err)
	}

	log.Debugf("got output from create user: %+v", out)

	return out.User, nil
}

func (i *IAM) DeleteUser(ctx context.Context, name string) error {
	if name == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if _, err := i.Service.DeleteUserWithContext(ctx, &iam.DeleteUserInput{
		UserName: aws.String(name),
	}); err != nil {
		return ErrCode("failed to delete user", err)
	}

	return nil
}

func (i *IAM) WaitForUser(ctx context.Context, name string) error {
	if name == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("waiting for user %s", name)

	if err := i.Service.WaitUntilUserExistsWithContext(ctx, &iam.GetUserInput{
		UserName: aws.String(name),
	}); err != nil {
		return ErrCode("failed waiting for user", err)
	}

	return nil
}

func (i *IAM) ListGroupsForUser(ctx context.Context, name string) ([]string, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("listing groups for user %s", name)

	out, err := i.Service.ListGroupsForUserWithContext(ctx, &iam.ListGroupsForUserInput{
		UserName: aws.String(name),
	})

	if err != nil {
		return nil, ErrCode("failed to list groups for user", err)
	}

	log.Debugf("got output listing groups for %s: %+v", name, out)

	groups := make([]string, 0, len(out.Groups))
	for _, g := range out.Groups {
		groups = append(groups, aws.StringValue(g.GroupName))
	}

	return groups, nil
}

func (i *IAM) TagUser(ctx context.Context, name string, tags []*iam.Tag) error {
	if name == "" || tags == nil {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("tagging user %s with tags %+v", name, tags)

	_, err := i.Service.TagUserWithContext(ctx, &iam.TagUserInput{
		UserName: aws.String(name),
		Tags:     tags,
	})

	if err != nil {
		return ErrCode("failed to tag user", err)
	}

	return nil
}
