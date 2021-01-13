package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/ecr"
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
