package ecr

import (
	"context"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	log "github.com/sirupsen/logrus"
)

// ListImages lists the images in a repostitory
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

// GetImages gets details about images in a repository
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

// GetImageScanFindings gets the scan findings for an image tag
func (e *ECR) GetImageScanFindings(ctx context.Context, repoName, tag string) (*ecr.ImageScanFindings, error) {
	if repoName == "" || tag == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting image scan findings for %s:%s", repoName, tag)

	out, err := e.Service.DescribeImageScanFindingsWithContext(ctx, &ecr.DescribeImageScanFindingsInput{
		ImageId: &ecr.ImageIdentifier{
			ImageTag: aws.String(tag),
		},
		MaxResults:     aws.Int64(1000),
		RepositoryName: aws.String(repoName),
	})

	if err != nil {
		return nil, ErrCode("failed to get image scan findings", err)
	}

	log.Debugf("got output from image scan findings %+v", out)

	return out.ImageScanFindings, nil
}

// DeleteImageTag deletes the image tag, if no other tags reference the image, the image is deleted
func (e *ECR) DeleteImageTag(ctx context.Context, repoName, tag string) (*ecr.BatchDeleteImageOutput, error) {
	if repoName == "" || tag == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("deleting image tag %s:%s", repoName, tag)

	out, err := e.Service.BatchDeleteImageWithContext(ctx, &ecr.BatchDeleteImageInput{
		ImageIds: []*ecr.ImageIdentifier{
			{
				ImageTag: aws.String(tag),
			},
		},
		RepositoryName: aws.String(repoName),
	})

	if err != nil {
		return nil, ErrCode("failed to delete image tag", err)
	}

	log.Debugf("got output from deleting image tag %+v", out)

	return out, nil
}
