package iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
)

var rootGroup01 = iam.Group{
	Arn:        aws.String("arn:aws:iam::12345678910:group/rootgroup01"),
	CreateDate: &testTime,
	GroupId:    aws.String("ROOTGROUP123"),
	GroupName:  aws.String("rootgroup01"),
	Path:       aws.String("/"),
}

var pathGroup01 = iam.Group{
	Arn:        aws.String("arn:aws:iam::12345678910:group/mypath/pathgroup01"),
	CreateDate: &testTime,
	GroupId:    aws.String("PATHGROUP123"),
	GroupName:  aws.String("pathgroup01"),
	Path:       aws.String("/mypath/"),
}

var testGroups = []*iam.Group{
	&rootGroup01,
	&pathGroup01,
}

func (m *mockIAMClient) GetGroupWithContext(ctx context.Context, input *iam.GetGroupInput, opts ...request.Option) (*iam.GetGroupOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, g := range testGroups {
		if aws.StringValue(input.GroupName) == aws.StringValue(g.GroupName) {
			return &iam.GetGroupOutput{Group: g}, nil
		}
	}

	return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "Not Found", nil)
}

func TestIAM_GetGroupWithPath(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.Group
		err     error
		wantErr bool
	}{
		{
			name: "empty name and path",
			args: args{
				ctx:  context.TODO(),
				name: "",
				path: "",
			},
			wantErr: true,
		},
		{
			name: "valid name, empty path",
			args: args{
				ctx:  context.TODO(),
				name: "rootgroup01",
				path: "",
			},
			want: &rootGroup01,
		},
		{
			name: "valid name, valid path",
			args: args{
				ctx:  context.TODO(),
				name: "pathgroup01",
				path: "/mypath/",
			},
			want: &pathGroup01,
		},
		{
			name: "missing group",
			args: args{
				ctx:  context.TODO(),
				name: "missing",
				path: "",
			},
			wantErr: true,
		},
		{
			name: "valid name, wrong path",
			args: args{
				ctx:  context.TODO(),
				name: "pathgroup01",
				path: "/wrongpath/",
			},
			wantErr: true,
		},
		{
			name: "api error",
			args: args{
				ctx:  context.TODO(),
				name: "rootgroup01",
				path: "",
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.GetGroupWithPath(tt.args.ctx, tt.args.name, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.GetGroupWithPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.GetGroupWithPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_CreateGroup(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.Group
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.CreateGroup(tt.args.ctx, tt.args.name, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.CreateGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_AttachGroupPolicy(t *testing.T) {
	type args struct {
		ctx       context.Context
		groupName string
		policyArn string
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
			if err := i.AttachGroupPolicy(tt.args.ctx, tt.args.groupName, tt.args.policyArn); (err != nil) != tt.wantErr {
				t.Errorf("IAM.AttachGroupPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIAM_ListAttachedGroupPolicies(t *testing.T) {
	type args struct {
		ctx       context.Context
		groupName string
		path      string
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
			got, err := i.ListAttachedGroupPolicies(tt.args.ctx, tt.args.groupName, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.ListAttachedGroupPolicies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.ListAttachedGroupPolicies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_AddUserToGroup(t *testing.T) {
	type args struct {
		ctx       context.Context
		userName  string
		groupName string
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
			if err := i.AddUserToGroup(tt.args.ctx, tt.args.userName, tt.args.groupName); (err != nil) != tt.wantErr {
				t.Errorf("IAM.AddUserToGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIAM_RemoveUserFromGroup(t *testing.T) {
	type args struct {
		ctx       context.Context
		userName  string
		groupName string
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
			if err := i.RemoveUserFromGroup(tt.args.ctx, tt.args.userName, tt.args.groupName); (err != nil) != tt.wantErr {
				t.Errorf("IAM.RemoveUserFromGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
