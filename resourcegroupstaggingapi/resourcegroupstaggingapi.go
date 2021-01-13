package resourcegroupstaggingapi

import (
	"context"
	"strings"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
	log "github.com/sirupsen/logrus"
)

// ResourceGroupsTaggingAPI is a wrapper around the aws resourcegroupstaggingapi service with some default config info
type ResourceGroupsTaggingAPI struct {
	session *session.Session
	Service resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
}

type ResourceGroupsTaggingAPIOption func(*ResourceGroupsTaggingAPI)

// Tag Filter is used to filter resources based on tags.  The Value portion is optional.
type TagFilter struct {
	Key   string
	Value []string
}

func New(opts ...ResourceGroupsTaggingAPIOption) ResourceGroupsTaggingAPI {
	r := ResourceGroupsTaggingAPI{}

	for _, opt := range opts {
		opt(&r)
	}

	if r.session != nil {
		r.Service = resourcegroupstaggingapi.New(r.session)
	}

	return r
}

func WithSession(sess *session.Session) ResourceGroupsTaggingAPIOption {
	return func(r *ResourceGroupsTaggingAPI) {
		log.Debug("using aws session")
		r.session = sess
	}
}

func WithCredentials(key, secret, token, region string) ResourceGroupsTaggingAPIOption {
	return func(r *ResourceGroupsTaggingAPI) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		r.session = sess
	}
}

// GetResourcesWithTags returns all of the resources with a type in the list of types that matches the tagfilters.  More
// details about which services support the resourgroup tagging api is here https://docs.aws.amazon.com/ARG/latest/userguide/supported-resources.html
func (r *ResourceGroupsTaggingAPI) GetResourcesWithTags(ctx context.Context, types []string, filters []*TagFilter) ([]*resourcegroupstaggingapi.ResourceTagMapping, error) {
	if len(filters) == 0 {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting resources with type '%s' that match tags", strings.Join(types, ", "))

	tagFilters := make([]*resourcegroupstaggingapi.TagFilter, 0, len(filters))
	for _, f := range filters {
		log.Debugf("tagfilter: %s:%+v", f.Key, f.Value)
		tagFilters = append(tagFilters, &resourcegroupstaggingapi.TagFilter{
			Key:    aws.String(f.Key),
			Values: aws.StringSlice(f.Value),
		})
	}

	out, err := r.Service.GetResourcesWithContext(ctx, &resourcegroupstaggingapi.GetResourcesInput{
		ResourceTypeFilters: aws.StringSlice(types),
		TagFilters:          tagFilters,
	})
	if err != nil {
		return nil, ErrCode("getting resource with tags", err)
	}

	log.Debugf("got output from get resources: %+v", out)

	return out.ResourceTagMappingList, nil
}
