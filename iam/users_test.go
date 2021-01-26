package iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestIAM_ListUsers(t *testing.T) {
	type args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.ListUsers(tt.args.ctx, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.ListUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_GetUserWithPath(t *testing.T) {
	type args struct {
		ctx  context.Context
		path string
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.User
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.GetUserWithPath(tt.args.ctx, tt.args.path, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.GetUserWithPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.GetUserWithPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_ListAccessKeys(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []*iam.AccessKeyMetadata
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.ListAccessKeys(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.ListAccessKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.ListAccessKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_CreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		path string
		tags []*iam.Tag
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.User
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.CreateUser(tt.args.ctx, tt.args.name, tt.args.path, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_DeleteUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
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
			if err := i.DeleteUser(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("IAM.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIAM_WaitForUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
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
			if err := i.WaitForUser(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("IAM.WaitForUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIAM_ListGroupsForUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.ListGroupsForUser(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.ListGroupsForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.ListGroupsForUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
