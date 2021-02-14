package iam

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

var testTime = time.Now()
var testPastTime = time.Unix(rand.Int63n(time.Now().Unix()), 0)

// mockIAMClient is a fake IAM client
type mockIAMClient struct {
	iamiface.IAMAPI
	t   *testing.T
	err error
}

func newMockIAMClient(t *testing.T, err error) iamiface.IAMAPI {
	return &mockIAMClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	client := New()
	to := reflect.TypeOf(client).String()
	if to != "iam.IAM" {
		t.Errorf("expected type to be iam.IAM, got %s", to)
	}
}
