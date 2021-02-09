package resourcegroupstaggingapi

import (
	"context"
	"reflect"
	"testing"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
	"github.com/pkg/errors"
)

// mockResourceGroupsTaggingAPIClient is a fake resourcegroupstaggingapi client
type mockResourceGroupsTaggingAPIClient struct {
	resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
	t   *testing.T
	err error
}

func newmockResourceGroupsTaggingAPIClient(t *testing.T, err error) resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI {
	return &mockResourceGroupsTaggingAPIClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	r := New()
	to := reflect.TypeOf(r).String()
	if to != "resourcegroupstaggingapi.ResourceGroupsTaggingAPI" {
		t.Errorf("expected type to be 'resourcegroupstaggingapi.ResourceGroupsTaggingAPI', got %s", to)
	}
}

type tag struct {
	key   string
	value string
}

type testResource struct {
	resourceType string
	tags         []tag
	arn          string
}

var testResources = []testResource{
	{
		resourceType: "ec2:instance",
		tags: []tag{
			{
				key:   "spinup:org",
				value: "foobar",
			},
			{
				key:   "spinup:spaceid",
				value: "123",
			},
		},
		arn: "arn:aws:ec2:us-east-1:1234567890:instance/i-0987654321",
	},
	{
		resourceType: "elasticloadbalancing:targetgroup",
		tags: []tag{
			{
				key:   "spinup:org",
				value: "foobar",
			},
			{
				key:   "spinup:spaceid",
				value: "123",
			},
		},
		arn: "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/testtg123/0987654321",
	},
	{
		resourceType: "elasticloadbalancing:targetgroup",
		tags: []tag{
			{
				key:   "spinup:org",
				value: "foobar",
			},
			{
				key:   "spinup:spaceid",
				value: "321",
			},
		},
		arn: "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/testtg321/0987654321",
	},
}

func (m *mockResourceGroupsTaggingAPIClient) GetResourcesWithContext(ctx context.Context, input *resourcegroupstaggingapi.GetResourcesInput, opts ...request.Option) (*resourcegroupstaggingapi.GetResourcesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	resourceList := []*resourcegroupstaggingapi.ResourceTagMapping{}
	for _, r := range testResources {
		if len(input.ResourceTypeFilters) > 0 {
			var typeMatch bool
			for _, t := range input.ResourceTypeFilters {
				if aws.StringValue(t) == r.resourceType {
					typeMatch = true
					break
				}
			}

			if !typeMatch {
				continue
			}
		}

		matches := true
		for _, filter := range input.TagFilters {
			innerMatch := func() bool {
				for _, rt := range r.tags {
					if aws.StringValue(filter.Key) == rt.key {
						if len(filter.Values) == 0 {
							return true
						}

						for _, value := range aws.StringValueSlice(filter.Values) {
							if value == rt.value {
								return true
							}
						}
					}
				}
				return false
			}()

			if !innerMatch {
				matches = false
			}
		}

		if matches {
			resourceList = append(resourceList, &resourcegroupstaggingapi.ResourceTagMapping{
				ResourceARN: aws.String(r.arn),
			})
		}
	}

	return &resourcegroupstaggingapi.GetResourcesOutput{
		ResourceTagMappingList: resourceList,
	}, nil
}

func TestGetResourcesWithTags(t *testing.T) {
	r := ResourceGroupsTaggingAPI{Service: newmockResourceGroupsTaggingAPIClient(t, nil)}
	filters := []*TagFilter{
		{
			Key:   "spinup:org",
			Value: []string{"foobar"},
		},
		{
			Key:   "spinup:spaceid",
			Value: []string{"123"},
		},
	}
	out, err := r.GetResourcesWithTags(context.TODO(), []string{"elasticloadbalancing:targetgroup"}, filters)
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	expected := []*resourcegroupstaggingapi.ResourceTagMapping{
		{
			ResourceARN: aws.String("arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/testtg123/0987654321"),
		},
	}

	if !reflect.DeepEqual(expected, out) {
		t.Errorf("expected %+v, got %+v", expected, out)
	}

	out, err = r.GetResourcesWithTags(context.TODO(), []string{}, filters)
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	expected = []*resourcegroupstaggingapi.ResourceTagMapping{
		{
			ResourceARN: aws.String("arn:aws:ec2:us-east-1:1234567890:instance/i-0987654321"),
		},
		{
			ResourceARN: aws.String("arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/testtg123/0987654321"),
		},
	}
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("expected %+v, got %+v", expected, out)
	}

	if _, err := r.GetResourcesWithTags(context.TODO(), []string{}, nil); err != nil {
		if aerr, ok := err.(apierror.Error); ok {
			if aerr.Code != apierror.ErrBadRequest {
				t.Errorf("expected error code %s, got: %s", apierror.ErrInternalError, aerr.Code)
			}
		} else {
			t.Errorf("expected apierror.Error")
		}
	} else {
		t.Error("expected error for empty filter list, got nil")
	}

	r.Service.(*mockResourceGroupsTaggingAPIClient).err = awserr.New(resourcegroupstaggingapi.ErrCodeThrottledException, "throttled", nil)
	if _, err := r.GetResourcesWithTags(context.TODO(), []string{}, filters); err != nil {
		if aerr, ok := err.(apierror.Error); ok {
			if aerr.Code != apierror.ErrConflict {
				t.Errorf("expected error code %s, got: %s", apierror.ErrConflict, aerr.Code)
			}
		} else {
			t.Errorf("expected apierror.Error")
		}
	} else {
		t.Error("expected error for empty filter list, got nil")
	}

	// test non-aws error
	r.Service.(*mockResourceGroupsTaggingAPIClient).err = errors.New("things blowing up!")
	if _, err := r.GetResourcesWithTags(context.TODO(), []string{}, filters); err != nil {
		if aerr, ok := err.(apierror.Error); ok {
			if aerr.Code != apierror.ErrInternalError {
				t.Errorf("expected error code %s, got: %s", apierror.ErrInternalError, aerr.Code)
			}
		} else {
			t.Errorf("expected apierror.Error")
		}
	} else {
		t.Error("expected error for empty filter list, got nil")
	}
}
