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

var rootuser1 = &iam.User{
	Arn:        aws.String("arn:aws:iam::0123456789:user/rootuser1"),
	CreateDate: aws.Time(testTime),
	Path:       aws.String("/"),
	UserId:     aws.String("ABCDEFGROOTUSER1"),
	UserName:   aws.String("rootuser1"),
}

var user1 = &iam.User{
	Arn:        aws.String("arn:aws:iam::0123456789:user/path1/user1"),
	CreateDate: aws.Time(testTime),
	Path:       aws.String("/path1/"),
	UserId:     aws.String("ABCDEFGUSER1"),
	UserName:   aws.String("user1"),
}

var user2 = &iam.User{
	Arn:        aws.String("arn:aws:iam::0123456789:user/path1/user2"),
	CreateDate: aws.Time(testTime),
	Path:       aws.String("/path1/"),
	UserId:     aws.String("ABCDEFGUSER2"),
	UserName:   aws.String("user2"),
}

var user3 = &iam.User{
	Arn:        aws.String("arn:aws:iam::0123456789:user/path1/user3"),
	CreateDate: aws.Time(testTime),
	Path:       aws.String("/path2/"),
	UserId:     aws.String("ABCDEFGUSER3"),
	UserName:   aws.String("user3"),
}

var testUsers = []*iam.User{
	rootuser1,
	user1,
	user2,
	user3,
}

var testAccessKeys = map[string][]*iam.AccessKeyMetadata{
	"rootuser1": {},
	"user1": {
		{
			AccessKeyId: aws.String("USER1XXXXXXXXX01"),
			CreateDate:  aws.Time(testPastTime),
			Status:      aws.String("Active"),
			UserName:    aws.String("user1"),
		},
		{
			AccessKeyId: aws.String("USER1XXXXXXXXX02"),
			CreateDate:  aws.Time(testPastTime),
			Status:      aws.String("Inactive"),
			UserName:    aws.String("user1"),
		},
	},
	"user2": {
		{
			AccessKeyId: aws.String("USER2XXXXXXXXX01"),
			CreateDate:  aws.Time(testPastTime),
			Status:      aws.String("Active"),
			UserName:    aws.String("user2"),
		},
		{
			AccessKeyId: aws.String("USER2XXXXXXXXX02"),
			CreateDate:  aws.Time(testPastTime),
			Status:      aws.String("Inactive"),
			UserName:    aws.String("user2"),
		},
	},
	"user3": {
		{
			AccessKeyId: aws.String("USER3XXXXXXXXX01"),
			CreateDate:  aws.Time(testPastTime),
			Status:      aws.String("InActive"),
			UserName:    aws.String("user3"),
		},
	},
}

var testUserGroups = map[string][]*iam.Group{
	"rootuser1": {
		{
			Arn:        aws.String(""),
			CreateDate: aws.Time(testTime),
			GroupId:    aws.String(""),
			GroupName:  aws.String("rootGroup1"),
			Path:       aws.String("/"),
		},
		{
			Arn:        aws.String(""),
			CreateDate: aws.Time(testTime),
			GroupId:    aws.String(""),
			GroupName:  aws.String("rootGroup2"),
			Path:       aws.String("/"),
		},
	},
	"user1": {
		{
			Arn:        aws.String(""),
			CreateDate: aws.Time(testTime),
			GroupId:    aws.String(""),
			GroupName:  aws.String("userGroup1"),
			Path:       aws.String("/path1/"),
		},
	},
	"user2": {
		{
			Arn:        aws.String(""),
			CreateDate: aws.Time(testTime),
			GroupId:    aws.String(""),
			GroupName:  aws.String("userGroup1"),
			Path:       aws.String("/path1/"),
		},
	},
	"user3": {
		{
			Arn:        aws.String(""),
			CreateDate: aws.Time(testTime),
			GroupId:    aws.String(""),
			GroupName:  aws.String("userGroup3"),
			Path:       aws.String("/path2/"),
		},
	},
}

func (m *mockIAMClient) ListUsersWithContext(ctx context.Context, input *iam.ListUsersInput, opts ...request.Option) (*iam.ListUsersOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	var users []*iam.User
	for _, u := range testUsers {
		if aws.StringValue(input.PathPrefix) == aws.StringValue(u.Path) {
			users = append(users, u)
		}
	}

	return &iam.ListUsersOutput{Users: users}, nil
}

func (m *mockIAMClient) GetUserWithContext(ctx context.Context, input *iam.GetUserInput, opts ...request.Option) (*iam.GetUserOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, u := range testUsers {
		if aws.StringValue(input.UserName) == aws.StringValue(u.UserName) {
			return &iam.GetUserOutput{User: u}, nil
		}
	}

	return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "Not Found", nil)
}

func (m *mockIAMClient) ListAccessKeysWithContext(ctx context.Context, input *iam.ListAccessKeysInput, opts ...request.Option) (*iam.ListAccessKeysOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for userName, keys := range testAccessKeys {
		if aws.StringValue(input.UserName) == userName {
			return &iam.ListAccessKeysOutput{AccessKeyMetadata: keys}, nil
		}
	}

	return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "Not Found", nil)
}

func (m *mockIAMClient) CreateUserWithContext(ctx context.Context, input *iam.CreateUserInput, opts ...request.Option) (*iam.CreateUserOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, u := range testUsers {
		iu := aws.StringValue(input.UserName)
		ou := aws.StringValue(u.UserName)
		ip := aws.StringValue(input.Path)
		op := aws.StringValue(u.Path)
		if (iu == ou) && (ip == op) {
			return &iam.CreateUserOutput{
				User: &iam.User{
					Arn:        u.Arn,
					CreateDate: u.CreateDate,
					Path:       u.Path,
					Tags:       input.Tags,
					UserId:     u.UserId,
					UserName:   u.UserName,
				},
			}, nil
		}
	}

	return &iam.CreateUserOutput{}, nil
}

func (m *mockIAMClient) DeleteUserWithContext(ctx context.Context, input *iam.DeleteUserInput, opts ...request.Option) (*iam.DeleteUserOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, u := range testUsers {
		if aws.StringValue(input.UserName) == aws.StringValue(u.UserName) {
			return &iam.DeleteUserOutput{}, nil
		}
	}

	return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "Not Found", nil)
}

func (m *mockIAMClient) ListGroupsForUserWithContext(ctx context.Context, input *iam.ListGroupsForUserInput, opts ...request.Option) (*iam.ListGroupsForUserOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for userName, groups := range testUserGroups {
		if aws.StringValue(input.UserName) == userName {
			return &iam.ListGroupsForUserOutput{Groups: groups}, nil
		}
	}

	return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "Not Found", nil)
}

func (m *mockIAMClient) DeleteAccessKeyWithContext(ctx context.Context, input *iam.DeleteAccessKeyInput, opts ...request.Option) (*iam.DeleteAccessKeyOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for userName, keys := range testAccessKeys {
		if aws.StringValue(input.UserName) == userName {
			for _, k := range keys {
				if aws.StringValue(k.AccessKeyId) == aws.StringValue(input.AccessKeyId) {
					if aws.StringValue(k.Status) != "Inactive" {
						return nil, awserr.New(iam.ErrCodeDeleteConflictException, "access key must be inactive", nil)
					}
					return &iam.DeleteAccessKeyOutput{}, nil
				}
			}
		}
	}

	return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "Not Found", nil)
}

func (m *mockIAMClient) CreateAccessKeyWithContext(ctx context.Context, input *iam.CreateAccessKeyInput, opts ...request.Option) (*iam.CreateAccessKeyOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, u := range testUsers {
		if aws.StringValue(input.UserName) == aws.StringValue(u.UserName) {
			return &iam.CreateAccessKeyOutput{
				AccessKey: &iam.AccessKey{
					CreateDate: aws.Time(testTime),
					UserName:   u.UserName,
					Status:     aws.String("Active"),
				},
			}, nil
		}
	}

	return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "Not Found", nil)
}

func (m *mockIAMClient) TagUserWithContext(ctx context.Context, input *iam.TagUserInput, opts ...request.Option) (*iam.TagUserOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, u := range testUsers {
		if aws.StringValue(input.UserName) == aws.StringValue(u.UserName) {
			return &iam.TagUserOutput{}, nil
		}
	}

	return nil, awserr.New(iam.ErrCodeNoSuchEntityException, "Not Found", nil)
}

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
		{
			name: "empty path",
			args: args{
				ctx:  context.TODO(),
				path: "",
			},
			want: []string{"rootuser1"},
		},
		{
			name: "root path",
			args: args{
				ctx:  context.TODO(),
				path: "/",
			},
			want: []string{"rootuser1"},
		},
		{
			name: "path1",
			args: args{
				ctx:  context.TODO(),
				path: "/path1/",
			},
			want: []string{"user1", "user2"},
		},
		{
			name: "path2",
			args: args{
				ctx:  context.TODO(),
				path: "/path2/",
			},
			want: []string{"user3"},
		},
		{
			name: "aws error",
			args: args{
				ctx:  context.TODO(),
				path: "/",
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
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
		{
			name: "empty path and name",
			args: args{
				ctx:  context.TODO(),
				path: "",
				name: "",
			},
			wantErr: true,
		},
		{
			name: "empty name",
			args: args{
				ctx:  context.TODO(),
				path: "/",
				name: "",
			},
			wantErr: true,
		},
		{
			name: "empty path, root user",
			args: args{
				ctx:  context.TODO(),
				path: "",
				name: "rootuser1",
			},
			want: rootuser1,
		},
		{
			name: "path1, user1",
			args: args{
				ctx:  context.TODO(),
				path: "/path1/",
				name: "user1",
			},
			want: user1,
		},
		{
			name: "path1, rootuser1",
			args: args{
				ctx:  context.TODO(),
				path: "/path1/",
				name: "rootuser1",
			},
			wantErr: true,
		},
		{
			name: "path2, user3",
			args: args{
				ctx:  context.TODO(),
				path: "/path2/",
				name: "user3",
			},
			want: user3,
		},
		{
			name: "aws error",
			args: args{
				ctx:  context.TODO(),
				path: "/",
				name: "rootuser1",
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
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
		{
			name: "empty name",
			args: args{
				ctx:  context.TODO(),
				name: "",
			},
			wantErr: true,
		},
		{
			name: "rootuser",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
			},
			want: testAccessKeys["rootuser1"],
		},
		{
			name: "user1",
			args: args{
				ctx:  context.TODO(),
				name: "user1",
			},
			want: testAccessKeys["user1"],
		},
		{
			name: "user2",
			args: args{
				ctx:  context.TODO(),
				name: "user2",
			},
			want: testAccessKeys["user2"],
		},
		{
			name: "user3",
			args: args{
				ctx:  context.TODO(),
				name: "user3",
			},
			want: testAccessKeys["user3"],
		},
		{
			name: "unknown user",
			args: args{
				ctx:  context.TODO(),
				name: "someotheruser",
			},
			wantErr: true,
		},
		{
			name: "aws error",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
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
		{
			name: "empty name, path and tags",
			args: args{
				name: "",
				path: "",
				tags: nil,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			args: args{
				name: "",
				path: "/path1/",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty path",
			args: args{
				name: "rootuser1",
				path: "",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			want: rootuser1,
		},
		{
			name: "empty tags",
			args: args{
				name: "rootuser1",
				path: "/",
				tags: nil,
			},
			want: rootuser1,
		},
		{
			name: "rootuser1 in /",
			args: args{
				name: "rootuser1",
				path: "/",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			want: rootuser1,
		},
		{
			name: "user1 in /path1/",
			args: args{
				name: "user1",
				path: "/path1/",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			want: user1,
		},
		{
			name: "user2 in /path1/",
			args: args{
				name: "user2",
				path: "/path1/",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			want: user2,
		},
		{
			name: "user3 in /path2/",
			args: args{
				name: "user3",
				path: "/path2/",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			want: user3,
		},
		{
			name: "aws error",
			args: args{
				name: "rootuser1",
				path: "/",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.CreateUser(tt.args.ctx, tt.args.name, tt.args.path, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// apply the tags passed with the args to the output (struct not pointer)
			var want *iam.User
			if tt.want != nil {
				w := *tt.want
				if tt.args.tags != nil {
					w.Tags = tt.args.tags
				}
				want = &w
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("IAM.CreateUser() = %v, want %v", got, want)
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
		{
			name: "empty name",
			args: args{
				ctx:  context.TODO(),
				name: "",
			},
			wantErr: true,
		},
		{
			name: "rootuser1",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
			},
		},
		{
			name: "user1",
			args: args{
				ctx:  context.TODO(),
				name: "user1",
			},
		},
		{
			name: "user2",
			args: args{
				ctx:  context.TODO(),
				name: "user2",
			},
		},
		{
			name: "user3",
			args: args{
				ctx:  context.TODO(),
				name: "user3",
			},
		},
		{
			name: "unknown user",
			args: args{
				ctx:  context.TODO(),
				name: "otheruser",
			},
			wantErr: true,
		},
		{
			name: "aws error",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
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
		{
			name: "empty name",
			args: args{
				ctx:  context.TODO(),
				name: "",
			},
			wantErr: true,
		},
		{
			name: "rootuser1",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
			},
			want: []string{"rootGroup1", "rootGroup2"},
		},
		{
			name: "user1",
			args: args{
				ctx:  context.TODO(),
				name: "user1",
			},
			want: []string{"userGroup1"},
		},
		{
			name: "user2",
			args: args{
				ctx:  context.TODO(),
				name: "user2",
			},
			want: []string{"userGroup1"},
		},
		{
			name: "user3",
			args: args{
				ctx:  context.TODO(),
				name: "user3",
			},
			want: []string{"userGroup3"},
		},
		{
			name: "unkown user",
			args: args{
				ctx:  context.TODO(),
				name: "someotheruser",
			},
			wantErr: true,
		},
		{
			name: "aws error",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
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

func TestIAM_CreateAccessKey(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		err     error
		want    *iam.AccessKey
		wantErr bool
	}{
		{
			name: "empty name",
			args: args{
				ctx:  context.TODO(),
				name: "",
			},
			wantErr: true,
		},
		{
			name: "rootuser1",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
			},
			want: &iam.AccessKey{
				CreateDate: aws.Time(testTime),
				UserName:   aws.String("rootuser1"),
				Status:     aws.String("Active"),
			},
		},
		{
			name: "user1",
			args: args{
				ctx:  context.TODO(),
				name: "user1",
			},
			want: &iam.AccessKey{
				CreateDate: aws.Time(testTime),
				UserName:   aws.String("user1"),
				Status:     aws.String("Active"),
			},
		},
		{
			name: "user2",
			args: args{
				ctx:  context.TODO(),
				name: "user2",
			},
			want: &iam.AccessKey{
				CreateDate: aws.Time(testTime),
				UserName:   aws.String("user2"),
				Status:     aws.String("Active"),
			},
		},
		{
			name: "user3",
			args: args{
				ctx:  context.TODO(),
				name: "user3",
			},
			want: &iam.AccessKey{
				CreateDate: aws.Time(testTime),
				UserName:   aws.String("user3"),
				Status:     aws.String("Active"),
			},
		},
		{
			name: "unknown user",
			args: args{
				ctx:  context.TODO(),
				name: "someotheruser",
			},
			wantErr: true,
		},
		{
			name: "aws error",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.CreateAccessKey(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.CreateAccessKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.CreateAccessKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_DeleteAccessKey(t *testing.T) {
	type args struct {
		ctx   context.Context
		name  string
		keyId string
	}
	tests := []struct {
		name    string
		args    args
		err     error
		wantErr bool
	}{
		{
			name: "empy name and key id",
			args: args{
				ctx:   context.TODO(),
				name:  "",
				keyId: "",
			},
			wantErr: true,
		},
		{
			name: "empy name",
			args: args{
				ctx:   context.TODO(),
				name:  "",
				keyId: "USER1XXXXXXXXX01",
			},
			wantErr: true,
		},
		{
			name: "empy key id",
			args: args{
				ctx:   context.TODO(),
				name:  "user1",
				keyId: "",
			},
			wantErr: true,
		},
		{
			name: "user1 active key USER1XXXXXXXXX01",
			args: args{
				ctx:   context.TODO(),
				name:  "user1",
				keyId: "USER1XXXXXXXXX01",
			},
			wantErr: true,
		},
		{
			name: "user1 inactive key USER1XXXXXXXXX02",
			args: args{
				ctx:   context.TODO(),
				name:  "user1",
				keyId: "USER1XXXXXXXXX02",
			},
		},
		{
			name: "user2 active key USER2XXXXXXXXX01",
			args: args{
				ctx:   context.TODO(),
				name:  "user2",
				keyId: "USER2XXXXXXXXX01",
			},
			wantErr: true,
		},
		{
			name: "user2 inactive key USER2XXXXXXXXX02",
			args: args{
				ctx:   context.TODO(),
				name:  "user2",
				keyId: "USER2XXXXXXXXX02",
			},
		},
		{
			name: "user3 active key USER3XXXXXXXXX01",
			args: args{
				ctx:   context.TODO(),
				name:  "user3",
				keyId: "USER3XXXXXXXXX01",
			},
			wantErr: true,
		},
		{
			name: "unknown user",
			args: args{
				ctx:   context.TODO(),
				name:  "someotheruser",
				keyId: "USER1XXXXXXXXX01",
			},
			wantErr: true,
		},
		{
			name: "unknown key",
			args: args{
				ctx:   context.TODO(),
				name:  "user1",
				keyId: "xxxxxmissingxxxxx",
			},
			wantErr: true,
		},
		{
			name: "aws error",
			args: args{
				ctx:   context.TODO(),
				name:  "user1",
				keyId: "USER1XXXXXXXXX02",
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			if err := i.DeleteAccessKey(tt.args.ctx, tt.args.name, tt.args.keyId); (err != nil) != tt.wantErr {
				t.Errorf("IAM.DeleteAccessKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIAM_TagUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		tags []*iam.Tag
	}
	tests := []struct {
		name    string
		args    args
		err     error
		wantErr bool
	}{
		{
			name: "empty name and tags",
			args: args{
				ctx:  context.TODO(),
				name: "",
				tags: nil,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			args: args{
				ctx:  context.TODO(),
				name: "",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty tags",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
				tags: nil,
			},
			wantErr: true,
		},
		{
			name: "rootuser1",
			args: args{
				ctx:  context.TODO(),
				name: "rootuser1",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
		},
		{
			name: "user1",
			args: args{
				ctx:  context.TODO(),
				name: "user1",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
		},
		{
			name: "user2",
			args: args{
				ctx:  context.TODO(),
				name: "user2",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
		},
		{
			name: "user3",
			args: args{
				ctx:  context.TODO(),
				name: "user3",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
		},
		{
			name: "unknown user",
			args: args{
				ctx:  context.TODO(),
				name: "someotheruser",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "aws error",
			args: args{
				ctx:  context.TODO(),
				name: "user1",
				tags: []*iam.Tag{
					{
						Key:   aws.String("foo"),
						Value: aws.String("bar"),
					},
				},
			},
			err:     awserr.New(iam.ErrCodeLimitExceededException, "limit exceeded", nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			if err := i.TagUser(tt.args.ctx, tt.args.name, tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("IAM.TagUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
