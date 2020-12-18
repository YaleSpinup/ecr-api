package ecr

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	log "github.com/sirupsen/logrus"
)

func (e *ECR) CreateRepository(ctx context.Context, input *ecr.CreateRepositoryInput) (*ecr.Repository, error) {
	if input == nil {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("creating repository %s", aws.StringValue(input.RepositoryName))

	out, err := e.Service.CreateRepositoryWithContext(ctx, input)
	if err != nil {
		return nil, ErrCode("failed to create repository", err)
	}

	log.Debugf("got create repostitory details %+v", out)

	return out.Repository, nil
}

func (e *ECR) ListRepositories(ctx context.Context) ([]string, error) {
	log.Info("listing all repositories")

	repos := []string{}
	err := e.Service.DescribeRepositoriesPagesWithContext(ctx,
		&ecr.DescribeRepositoriesInput{MaxResults: aws.Int64(1000)},
		func(page *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
			for _, r := range page.Repositories {
				repos = append(repos, aws.StringValue(r.RepositoryName))
			}

			return true
		})

	log.Debugf("got list of repostitories %+v", repos)

	return repos, err
}

func (e *ECR) GetRepositories(ctx context.Context, repoName string) (*ecr.Repository, error) {
	if repoName == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting details for repostories: %s", repoName)

	out, err := e.Service.DescribeRepositoriesWithContext(ctx, &ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{
			aws.String(repoName),
		},
	})

	if err != nil {
		return nil, ErrCode("failed to get repository details", err)
	}

	log.Debugf("got repostitory details %+v", out)

	if len(out.Repositories) == 0 {
		msg := fmt.Sprintf("%s not found", repoName)
		return nil, apierror.New(apierror.ErrNotFound, msg, nil)
	}

	if num := len(out.Repositories); num > 1 {
		msg := fmt.Sprintf("unexpected number of repositories found for id %s (%d)", repoName, num)
		return nil, apierror.New(apierror.ErrInternalError, msg, nil)
	}

	return out.Repositories[0], nil
}

func (e *ECR) DeleteRepository(ctx context.Context, repoName string) (*ecr.Repository, error) {
	if repoName == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	out, err := e.Service.DeleteRepositoryWithContext(ctx, &ecr.DeleteRepositoryInput{
		Force:          aws.Bool(true),
		RepositoryName: aws.String(repoName),
	})

	if err != nil {
		return nil, ErrCode("failed to delete repository", err)
	}

	log.Debugf("got output from repository delete: %+v", out)

	return out.Repository, nil
}

func (e *ECR) GetRepositoryTags(ctx context.Context, repoArn string) ([]*ecr.Tag, error) {
	if repoArn == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting tags for repository %s", repoArn)

	out, err := e.Service.ListTagsForResourceWithContext(ctx, &ecr.ListTagsForResourceInput{
		ResourceArn: aws.String(repoArn),
	})

	if err != nil {
		return nil, ErrCode("failed to get repository tags", err)
	}

	log.Debugf("got repostitory tags %+v", out)

	return out.Tags, nil
}

func (e *ECR) UpdateRepositoryTags(ctx context.Context, repoArn string, tags []*ecr.Tag) error {
	if repoArn == "" || tags == nil {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("updating tags for repository %s", repoArn)

	out, err := e.Service.TagResourceWithContext(ctx, &ecr.TagResourceInput{
		ResourceArn: aws.String(repoArn),
		Tags:        tags,
	})
	if err != nil {
		return ErrCode("failed to update repository tags", err)
	}

	log.Debugf("got output from updating repostiory tags %+v", out)

	return nil
}

func (e *ECR) SetImageScanningConfiguration(ctx context.Context, input *ecr.PutImageScanningConfigurationInput) error {
	if input == nil {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("updating scanning configuration for repository %s", aws.StringValue(input.RepositoryName))

	out, err := e.Service.PutImageScanningConfigurationWithContext(ctx, input)
	if err != nil {
		return ErrCode("failed to update repository", err)
	}

	log.Debugf("got output from updating image scanning configuration %+v", out)

	return nil
}
