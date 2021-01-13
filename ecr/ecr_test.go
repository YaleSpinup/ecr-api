package ecr

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

// mockECRClient is a fake ecs client
type mockECRClient struct {
	ecriface.ECRAPI
	t   *testing.T
	err error
}

func newmockECRClient(t *testing.T, err error) ecriface.ECRAPI {
	return &mockECRClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	client := New()
	to := reflect.TypeOf(client).String()
	if to != "ecr.ECR" {
		t.Errorf("expected type to be 'ecs.ECS', got %s", to)
	}
}
