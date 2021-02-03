package iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestIAM_GetPolicyByName(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.Policy
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.GetPolicyByName(tt.args.ctx, tt.args.name, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.GetPolicyByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.GetPolicyByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_GetDefaultPolicyVersion(t *testing.T) {
	type args struct {
		ctx     context.Context
		arn     string
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.PolicyVersion
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.GetDefaultPolicyVersion(tt.args.ctx, tt.args.arn, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.GetDefaultPolicyVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.GetDefaultPolicyVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_WaitForPolicy(t *testing.T) {

	type args struct {
		ctx       context.Context
		policyArn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			if err := i.WaitForPolicy(tt.args.ctx, tt.args.policyArn); (err != nil) != tt.wantErr {
				t.Errorf("IAM.WaitForPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIAM_CreatePolicy(t *testing.T) {
	type args struct {
		ctx       context.Context
		name      string
		path      string
		policyDoc string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.Policy
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.CreatePolicy(tt.args.ctx, tt.args.name, tt.args.path, tt.args.policyDoc)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.CreatePolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.CreatePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_UpdatePolicy(t *testing.T) {
	type args struct {
		ctx       context.Context
		arn       string
		policyDoc string
	}
	tests := []struct {
		name    string
		args    args
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			if err := i.UpdatePolicy(tt.args.ctx, tt.args.arn, tt.args.policyDoc); (err != nil) != tt.wantErr {
				t.Errorf("IAM.UpdatePolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
