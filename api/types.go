package api

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/ecr"
	log "github.com/sirupsen/logrus"
)

type RepositoryCreateRequest struct {
	// Specify a custom KmsKeyId.  This will also change the encryption type from
	// 'AES256' to 'KMS'.  By default, when no encryption configuration is set or
	// the AES256 encryption type is used, Amazon ECR uses server-side encryption
	// with Amazon S3-managed encryption keys which encrypts your data at rest using
	// an AES-256 encryption algorithm.
	// Note: AWS KMS enforces a limit of 500 grants per CMK. As a result, there is
	// a limit of 500 Amazon ECR repositories that can be encrypted per CMK.
	KmsKeyId string

	// TODO Lifecycle Policy
	// TODO Repository Policy

	// The setting that determines whether images are scanned after being pushed
	// to a repository. If set to true, images will be scanned after being pushed.
	// If this parameter is not specified, it will default to false and images will
	// not be scanned unless a scan is manually started with the StartImageScan
	// API.
	ScanOnPush string

	// The name to use for the repository. The repository name may be specified
	// on its own (such as nginx-web-app) or it can be prepended with a namespace
	// to group the repository into a category (such as project-a/nginx-web-app)
	RepositoryName string

	// Tags to apply to the repository
	Tags []*Tag
}

type RepositoryUpdateRequest struct {
	ScanOnPush string
	Tags       []*Tag
}

type RepositoryResponse struct {
	CreatedAt          time.Time
	EncryptionType     string
	KmsKeyId           string
	ScanOnPush         string
	ImageTagMutability string
	RegistryId         string
	RepositoryArn      string
	RepositoryName     string
	RepositoryUri      string
	Tags               []*Tag
}

type Tag struct {
	Key   string
	Value string
}

// repositoryResponseFromECR maps ECR response to a common struct
func repositoryResponseFromECR(r *ecr.Repository, t []*ecr.Tag) *RepositoryResponse {
	log.Debugf("mapping repository %s", awsutil.Prettify(r))

	repository := RepositoryResponse{
		CreatedAt:          aws.TimeValue(r.CreatedAt),
		ImageTagMutability: aws.StringValue(r.ImageTagMutability),
		RegistryId:         aws.StringValue(r.RegistryId),
		RepositoryArn:      aws.StringValue(r.RepositoryArn),
		RepositoryName:     aws.StringValue(r.RepositoryName),
		RepositoryUri:      aws.StringValue(r.RepositoryUri),
		Tags:               fromECRTags(t),
	}

	if r.ImageScanningConfiguration != nil {
		b := aws.BoolValue(r.ImageScanningConfiguration.ScanOnPush)
		repository.ScanOnPush = strconv.FormatBool(b)
	}

	if r.EncryptionConfiguration != nil {
		repository.EncryptionType = aws.StringValue(r.EncryptionConfiguration.EncryptionType)
		repository.KmsKeyId = aws.StringValue(r.EncryptionConfiguration.KmsKey)
	}

	return &repository
}

// normalizTags strips the org, spaceid and name from the given tags and ensures they
// are set to the API org and the group string, name passed to the request
func normalizeTags(org, group, name string, tags []*Tag) []*Tag {
	normalizedTags := []*Tag{}
	for _, t := range tags {
		if t.Key == "spinup:spaceid" || t.Key == "spinup:org" || t.Key == "Name" {
			continue
		}
		normalizedTags = append(normalizedTags, t)
	}

	normalizedTags = append(normalizedTags,
		&Tag{
			Key:   "Name",
			Value: name,
		},
		&Tag{
			Key:   "spinup:org",
			Value: org,
		}, &Tag{
			Key:   "spinup:spaceid",
			Value: group,
		})

	log.Debugf("returning normalized tags: %+v", normalizedTags)
	return normalizedTags
}

// fromECRTags converts from ECR tags to api Tags
func fromECRTags(ecrTags []*ecr.Tag) []*Tag {
	tags := make([]*Tag, 0, len(ecrTags))
	for _, t := range ecrTags {
		tags = append(tags, &Tag{
			Key:   aws.StringValue(t.Key),
			Value: aws.StringValue(t.Value),
		})
	}
	return tags
}

// toECRTags converts from api Tags to ECR tags
func toECRTags(tags []*Tag) []*ecr.Tag {
	efsTags := make([]*ecr.Tag, 0, len(tags))
	for _, t := range tags {
		efsTags = append(efsTags, &ecr.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}
	return efsTags
}