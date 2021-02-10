package sts

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

var testTime = time.Now()

// mockSTSClient is a fake sts client
type mockSTSClient struct {
	stsiface.STSAPI
	t   *testing.T
	err error
}

func newMockSTSClient(t *testing.T, err error) stsiface.STSAPI {
	return &mockSTSClient{
		t:   t,
		err: err,
	}
}

var testAssumeRoleOutput = &sts.AssumeRoleOutput{
	AssumedRoleUser: &sts.AssumedRoleUser{
		Arn:           aws.String("arn:aws:sts::0123456789:assumed-role/UnitTestXAManagementRole/spinup-unit-ecr-api-000000-11111-2222-3333-444444"),
		AssumedRoleId: aws.String("AABBCCDDEEFFGGHHIIJJ12345:spinup-unit-ecr-api-000000-11111-2222-3333-444444"),
	},
	Credentials: &sts.Credentials{
		AccessKeyId:     aws.String(""),
		Expiration:      aws.Time(testTime),
		SecretAccessKey: aws.String(""),
		SessionToken:    aws.String(""),
	},
}

func (m *mockSTSClient) AssumeRoleWithContext(ctx context.Context, input *sts.AssumeRoleInput, opts ...request.Option) (*sts.AssumeRoleOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	return testAssumeRoleOutput, nil
}

func TestNew(t *testing.T) {
	client := New()
	to := reflect.TypeOf(client).String()
	if to != "sts.STS" {
		t.Errorf("expected type to be 'sts.STS', got %s", to)
	}
}

func TestSTS_AssumeRole(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *sts.AssumeRoleInput
	}
	tests := []struct {
		name    string
		args    args
		want    *sts.AssumeRoleOutput
		err     error
		wantErr bool
	}{
		{
			name: "nil input",
			args: args{
				ctx:   context.TODO(),
				input: nil,
			},
			wantErr: true,
		},
		{
			name: "empty role arn",
			args: args{
				ctx:   context.TODO(),
				input: &sts.AssumeRoleInput{},
			},
			wantErr: true,
		},
		{
			name: "valid role arn",
			args: args{
				ctx: context.TODO(),
				input: &sts.AssumeRoleInput{
					RoleArn: aws.String("arn:aws:iam::516855177326:role/UnitTestXAManagementRole"),
				},
			},
			want: testAssumeRoleOutput,
		},
		{
			name: "aws error",
			args: args{
				ctx: context.TODO(),
				input: &sts.AssumeRoleInput{
					RoleArn: aws.String("arn:aws:iam::516855177326:role/UnitTestXAManagementRole"),
				},
			},
			err:     awserr.New(sts.ErrCodeMalformedPolicyDocumentException, "bad policy yo", nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &STS{
				Service: newMockSTSClient(t, tt.err),
				Org:     "unit",
			}
			got, err := s.AssumeRole(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("STS.AssumeRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STS.AssumeRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
