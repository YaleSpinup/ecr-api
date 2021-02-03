package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/iam"

	log "github.com/sirupsen/logrus"
)

// RepositoryCreateRequest is the payload for creating an ECR repository
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

// RepositoryUpdateRequest is the payload for updating an ECR repository
type RepositoryUpdateRequest struct {
	ScanOnPush string
	Tags       []*Tag
}

// RepositoryResponse is the response payload for repository operations
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

// RepositoryUserCreateRequest is the request payload for creating a repository user
type RepositoryUserCreateRequest struct {
	UserName string
	Groups   []string
	Tags     []*Tag
}

// RepositoryUserResponse is the response payload for user operations
type RepositoryUserResponse struct {
	UserName   string
	AccessKeys []*iam.AccessKeyMetadata
	Groups     []string
	Tags       []*Tag
}

// Tag is our AWS compatible tag struct that can be converted to specific tag types
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

// repositoryUserResponseFromIAM maps IAM response to a common struct
func repositoryUserResponseFromIAM(org string, u *iam.User, keys []*iam.AccessKeyMetadata, groups []string) *RepositoryUserResponse {
	log.Debugf("mapping iam user %s", awsutil.Prettify(u))

	userName := aws.StringValue(u.UserName)

	// path is format: /spinup/%s/%s/%s/
	path := strings.Split(aws.StringValue(u.Path), "/")

	if len(path) > 2 {
		prefix := fmt.Sprintf("%s-%s-", path[len(path)-3], path[len(path)-2])

		log.Debugf("trimming prefix '%s' from username %s", prefix, userName)

		userName = strings.TrimPrefix(userName, prefix)
	}

	if keys == nil {
		keys = []*iam.AccessKeyMetadata{}
	}

	for i, g := range groups {
		groups[i] = strings.TrimSuffix(g, "-"+org)
	}

	user := RepositoryUserResponse{
		AccessKeys: keys,
		Groups:     groups,
		Tags:       fromIAMTags(u.Tags),
		UserName:   userName,
	}

	return &user
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
		},
		&Tag{
			Key:   "spinup:spaceid",
			Value: group,
		})

	log.Debugf("returning normalized tags: %+v", normalizedTags)
	return normalizedTags
}

// normalizUserTags strips the org, spaceid, resource, and name from the given tags
// and ensures they are set to the API org, group string, managed resource and name
// passed to the request
func normalizeUserTags(org, group, resource, name string, tags []*Tag) []*Tag {
	normalizedTags := []*Tag{}
	for _, t := range tags {
		if t.Key == "spinup:spaceid" || t.Key == "spinup:org" || t.Key == "ResourceName" || t.Key == "Name" {
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
			Key:   "ResourceName",
			Value: resource,
		},
		&Tag{
			Key:   "spinup:org",
			Value: org,
		},
		&Tag{
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
	ecrTags := make([]*ecr.Tag, 0, len(tags))
	for _, t := range tags {
		ecrTags = append(ecrTags, &ecr.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}
	return ecrTags
}

// fromIAMTags converts from IAM tags to api Tags
func fromIAMTags(iamTags []*iam.Tag) []*Tag {
	tags := make([]*Tag, 0, len(iamTags))
	for _, t := range iamTags {
		tags = append(tags, &Tag{
			Key:   aws.StringValue(t.Key),
			Value: aws.StringValue(t.Value),
		})
	}
	return tags
}

// toIAMTags converts from api Tags to IAM tags
func toIAMTags(tags []*Tag) []*iam.Tag {
	iamTags := make([]*iam.Tag, 0, len(tags))
	for _, t := range tags {
		iamTags = append(iamTags, &iam.Tag{
			Key:   aws.String(t.Key),
			Value: aws.String(t.Value),
		})
	}
	return iamTags
}
