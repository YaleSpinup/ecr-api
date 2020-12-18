package ecr

import (
	"context"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	log "github.com/sirupsen/logrus"
)

func (e *ECR) ListImages(ctx context.Context, repoName string) ([]*ecr.ImageIdentifier, error) {
	if repoName == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("listing images for repository %s", repoName)

	out, err := e.Service.ListImagesWithContext(ctx, &ecr.ListImagesInput{
		RepositoryName: aws.String(repoName),
	})
	if err != nil {
		return nil, ErrCode("failed to list repository images", err)
	}

	log.Debugf("got output from listing repostiory images %+v", out)

	return out.ImageIds, nil
}

func (e *ECR) GetImages(ctx context.Context, repoName string, imageIds ...*ecr.ImageIdentifier) ([]*ecr.ImageDetail, error) {
	if repoName == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("listing images for repository %s", repoName)

	input := &ecr.DescribeImagesInput{RepositoryName: aws.String(repoName)}
	if len(imageIds) > 0 {
		input.SetImageIds(imageIds)
	}

	out, err := e.Service.DescribeImagesWithContext(ctx, input)
	if err != nil {
		return nil, ErrCode("failed to get images", err)
	}

	log.Debugf("got output from images %+v", out)

	return out.ImageDetails, nil
}
