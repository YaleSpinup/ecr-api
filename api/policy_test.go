package api

import (
	"reflect"
	"testing"
)

func Test_orgTagAccessPolicy(t *testing.T) {
	type args struct {
		org string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "org policy",
			args: args{
				org: "testOrg",
			},
			want: `{"Version":"2012-10-17","Statement":[{"Action":["*"],"Condition":{"StringEquals":{"aws:ResourceTag/spinup:org":["testOrg"]}},"Effect":"Allow","Resource":["*"]}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orgTagAccessPolicy(tt.args.org)
			if (err != nil) != tt.wantErr {
				t.Errorf("orgTagAccessPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("orgTagAccessPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_repositoryUserCreatePolicy(t *testing.T) {
	type fields struct {
		org string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "test org",
			fields: fields{
				org: "testOrg",
			},
			want: `{"Version":"2012-10-17","Statement":[{"Action":["iam:CreatePolicy","iam:UntagUser","iam:GetPolicyVersion","iam:AddUserToGroup","iam:GetPolicy","iam:ListAttachedGroupPolicies","iam:ListGroupPolicies","iam:AttachGroupPolicy","iam:GetUser","iam:CreatePolicyVersion","iam:CreateUser","iam:GetGroup","iam:TagUser"],"Effect":"Allow","Resource":["arn:aws:iam::*:group/spinup/testOrg/*","arn:aws:iam::*:policy/spinup/testOrg/*","arn:aws:iam::*:user/spinup/testOrg/*"],"Sid":"CreateRepositoryUser"},{"Action":["iam:ListPolicies"],"Effect":"Allow","Resource":["*"],"Sid":"ListRepositoryUserPolicies"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				org: tt.fields.org,
			}
			got, err := s.repositoryUserCreatePolicy()
			if (err != nil) != tt.wantErr {
				t.Errorf("server.repositoryUserCreatePolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("server.repositoryUserCreatePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_repositoryUserDeletePolicy(t *testing.T) {
	type fields struct {
		org string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "test org",
			fields: fields{
				org: "testOrg",
			},
			want: `{"Version":"2012-10-17","Statement":[{"Action":["iam:DeleteAccessKey","iam:RemoveUserFromGroup","iam:ListAccessKeys","iam:ListGroupsForUser","iam:DeleteUser","iam:GetUser"],"Effect":"Allow","Resource":["arn:aws:iam::*:user/spinup/testOrg/*","arn:aws:iam::*:group/spinup/testOrg/SpinupECRAdminGroup-testOrg"],"Sid":"DeleteRepositoryUser"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				org: tt.fields.org,
			}
			got, err := s.repositoryUserDeletePolicy()
			if (err != nil) != tt.wantErr {
				t.Errorf("server.repositoryUserDeletePolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("server.repositoryUserDeletePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_repositoryUserUpdatePolicy(t *testing.T) {
	type fields struct {
		org string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "test org",
			fields: fields{
				org: "testOrg",
			},
			want: `{"Version":"2012-10-17","Statement":[{"Action":["iam:UntagUser","iam:DeleteAccessKey","iam:RemoveUserFromGroup","iam:TagUser","iam:CreateAccessKey","iam:ListAccessKeys"],"Effect":"Allow","Resource":["arn:aws:iam::*:user/spinup/testOrg/*","arn:aws:iam::*:group/spinup/testOrg/SpinupECRAdminGroup-testOrg"],"Sid":"UpdateRepositoryUser"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				org: tt.fields.org,
			}
			got, err := s.repositoryUserUpdatePolicy()
			if (err != nil) != tt.wantErr {
				t.Errorf("server.repositoryUserUpdatePolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("server.repositoryUserUpdatePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repositoryPolicy(t *testing.T) {
	type args struct {
		groups []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				groups: nil,
			},
			want: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":["${aws:ResourceTag/spinup:spaceid}"]}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
		},
		{
			name: "empty list",
			args: args{
				groups: []string{},
			},
			want: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":["${aws:ResourceTag/spinup:spaceid}"]}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
		},
		{
			name: "multiple groups",
			args: args{
				groups: []string{"foo", "bar", "baz"},
			},
			want: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":["${aws:ResourceTag/spinup:spaceid}","foo","bar","baz"]}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repositoryPolicy(tt.args.groups)
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repositoryPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repositoryGroupsFromPolicy(t *testing.T) {
	type args struct {
		policy string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				policy: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":["${aws:ResourceTag/spinup:spaceid}"]}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
			},
			want: []string{},
		},
		{
			name: "empty policy",
			args: args{
				policy: "",
			},
			want: []string{},
		},
		{
			name: "empty list",
			args: args{
				policy: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":["${aws:ResourceTag/spinup:spaceid}"]}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
			},
			want: []string{},
		},
		{
			name: "non-slice in tag list",
			args: args{
				policy: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":"somenotslice"}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
			},
			want: []string{},
		},
		{
			name: "non-strings slice in tag list",
			args: args{
				policy: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":[1,2,3]}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
			},
			want: []string{},
		},
		{
			name: "multiple groups",
			args: args{
				policy: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":["${aws:ResourceTag/spinup:spaceid}","foo","bar","baz"]}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
			},
			want: []string{"foo", "bar", "baz"},
		},
		{
			name: "unexpected policy SID",
			args: args{
				policy: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"StringEqualsIgnoreCase":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":["${aws:ResourceTag/spinup:spaceid}","foo","bar","baz"]}},"Effect":"Allow","Principal":"*","Sid":"SomeOtherSID"}]}`,
			},
			want: []string{},
		},
		{
			name: "missing StringEqualsIgnoreCase",
			args: args{
				policy: `{"Version":"2012-10-17","Statement":[{"Action":["ecr:GetAuthorizationToken","ecr:BatchCheckLayerAvailability","ecr:GetDownloadUrlForLayer","ecr:BatchGetImage"],"Condition":{"OtherFooCondition":{"aws:PrincipalTag/spinup:org":"${aws:ResourceTag/spinup:org}","aws:PrincipalTag/spinup:spaceid":["${aws:ResourceTag/spinup:spaceid}","foo","bar","baz"]}},"Effect":"Allow","Principal":"*","Sid":"AllowPullImagesFromSpaceAndOrg"}]}`,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repositoryGroupsFromPolicy(tt.args.policy)
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryGroupsFromPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repositoryGroupsFromPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
