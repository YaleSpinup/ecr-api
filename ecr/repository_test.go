package ecr

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/pkg/errors"
)

var tRepos = []*ecr.Repository{
	{
		EncryptionConfiguration: &ecr.EncryptionConfiguration{
			EncryptionType: aws.String("AES256"),
		},
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true),
		},
		ImageTagMutability: aws.String("MUTABLE"),
		RegistryId:         aws.String("012345678910"),
		RepositoryArn:      aws.String("arn:aws:ecr:us-east-1:012345678910:repository/carols/12DaysOfChristmas"),
		RepositoryName:     aws.String("carols/12DaysOfChristmas"),
		RepositoryUri:      aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/carols/12DaysOfChristmas"),
	},
	{
		EncryptionConfiguration: &ecr.EncryptionConfiguration{
			EncryptionType: aws.String("AES256"),
		},
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true),
		},
		ImageTagMutability: aws.String("MUTABLE"),
		RegistryId:         aws.String("012345678910"),
		RepositoryArn:      aws.String("arn:aws:ecr:us-east-1:012345678910:repository/carols/SilentNight"),
		RepositoryName:     aws.String("carols/SilentNight"),
		RepositoryUri:      aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/carols/SilentNight"),
	},
	{
		EncryptionConfiguration: &ecr.EncryptionConfiguration{
			EncryptionType: aws.String("AES256"),
		},
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true),
		},
		ImageTagMutability: aws.String("MUTABLE"),
		RegistryId:         aws.String("012345678910"),
		RepositoryArn:      aws.String("arn:aws:ecr:us-east-1:012345678910:repository/carols/FrostyTheSnowman"),
		RepositoryName:     aws.String("carols/FrostyTheSnowman"),
		RepositoryUri:      aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/carols/FrostyTheSnowman"),
	},
	{
		EncryptionConfiguration: &ecr.EncryptionConfiguration{
			EncryptionType: aws.String("AES256"),
		},
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true),
		},
		ImageTagMutability: aws.String("MUTABLE"),
		RegistryId:         aws.String("012345678910"),
		RepositoryArn:      aws.String("arn:aws:ecr:us-east-1:012345678910:repository/carols/LittleDrummerBoy"),
		RepositoryName:     aws.String("carols/LittleDrummerBoy"),
		RepositoryUri:      aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/carols/LittleDrummerBoy"),
	},
	{
		EncryptionConfiguration: &ecr.EncryptionConfiguration{
			EncryptionType: aws.String("AES256"),
		},
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true),
		},
		ImageTagMutability: aws.String("MUTABLE"),
		RegistryId:         aws.String("012345678910"),
		RepositoryArn:      aws.String("arn:aws:ecr:us-east-1:012345678910:repository/reindeer/rudolph"),
		RepositoryName:     aws.String("reindeer/rudolph"),
		RepositoryUri:      aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/reindeer/rudolph"),
	},
	{
		EncryptionConfiguration: &ecr.EncryptionConfiguration{
			EncryptionType: aws.String("AES256"),
		},
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true),
		},
		ImageTagMutability: aws.String("MUTABLE"),
		RegistryId:         aws.String("012345678910"),
		RepositoryArn:      aws.String("arn:aws:ecr:us-east-1:012345678910:repository/reindeer/dasher"),
		RepositoryName:     aws.String("reindeer/dasher"),
		RepositoryUri:      aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/reindeer/dasher"),
	},
	{
		EncryptionConfiguration: &ecr.EncryptionConfiguration{
			EncryptionType: aws.String("AES256"),
		},
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true),
		},
		ImageTagMutability: aws.String("MUTABLE"),
		RegistryId:         aws.String("012345678910"),
		RepositoryArn:      aws.String("arn:aws:ecr:us-east-1:012345678910:repository/reindeer/dancer"),
		RepositoryName:     aws.String("reindeer/dancer"),
		RepositoryUri:      aws.String("012345678910.dkr.ecr.us-east-1.amazonaws.com/reindeer/dancer"),
	},
}

func (m *mockECRClient) CreateRepositoryWithContext(ctx context.Context, input *ecr.CreateRepositoryInput, opts ...request.Option) (*ecr.CreateRepositoryOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	return &ecr.CreateRepositoryOutput{
		Repository: &ecr.Repository{
			EncryptionConfiguration:    input.EncryptionConfiguration,
			ImageScanningConfiguration: input.ImageScanningConfiguration,
			ImageTagMutability:         input.ImageTagMutability,
			RegistryId:                 aws.String("012345678910"),
			RepositoryArn:              aws.String(fmt.Sprintf("arn:aws:ecr:us-east-1:012345678910:repository/%s", aws.StringValue(input.RepositoryName))),
			RepositoryName:             input.RepositoryName,
			RepositoryUri:              aws.String(fmt.Sprintf("012345678910.dkr.ecr.us-east-1.amazonaws.com/%s", aws.StringValue(input.RepositoryName))),
		},
	}, nil
}

func (m *mockECRClient) DescribeRepositoriesPagesWithContext(ctx context.Context, input *ecr.DescribeRepositoriesInput, f func(*ecr.DescribeRepositoriesOutput, bool) bool, opts ...request.Option) error {
	if m.err != nil {
		return m.err
	}

	_ = f(&ecr.DescribeRepositoriesOutput{Repositories: tRepos}, false)

	return nil
}

func (m *mockECRClient) DescribeRepositoriesWithContext(ctx context.Context, input *ecr.DescribeRepositoriesInput, opts ...request.Option) (*ecr.DescribeRepositoriesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	// special case to return more than one repo
	if len(input.RepositoryNames) == 1 && aws.StringValue(input.RepositoryNames[0]) == "manyreposmatch" {
		return &ecr.DescribeRepositoriesOutput{Repositories: tRepos}, nil
	}

	repos := []*ecr.Repository{}
	for _, r := range tRepos {
		for _, i := range input.RepositoryNames {
			if aws.StringValue(i) == aws.StringValue(r.RepositoryName) {
				repos = append(repos, r)
			}
		}
	}

	return &ecr.DescribeRepositoriesOutput{Repositories: repos}, nil
}

func (m *mockECRClient) DeleteRepositoryWithContext(ctx context.Context, input *ecr.DeleteRepositoryInput, opts ...request.Option) (*ecr.DeleteRepositoryOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	if !aws.BoolValue(input.Force) {
		return nil, awserr.New(ecr.ErrCodeRepositoryNotEmptyException, "repository not empty", nil)
	}

	for _, r := range tRepos {
		if aws.StringValue(input.RepositoryName) == aws.StringValue(r.RepositoryName) {
			return &ecr.DeleteRepositoryOutput{}, nil
		}
	}

	return nil, awserr.New(ecr.ErrCodeRepositoryNotFoundException, "repository not found", nil)
}

func (m *mockECRClient) PutImageScanningConfigurationWithContext(ctx context.Context, input *ecr.PutImageScanningConfigurationInput, opts ...request.Option) (*ecr.PutImageScanningConfigurationOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, r := range tRepos {
		if aws.StringValue(input.RepositoryName) == aws.StringValue(r.RepositoryName) {
			return &ecr.PutImageScanningConfigurationOutput{
				ImageScanningConfiguration: input.ImageScanningConfiguration,
				RegistryId:                 r.RegistryId,
				RepositoryName:             r.RepositoryName,
			}, nil
		}
	}

	return nil, awserr.New(ecr.ErrCodeRepositoryNotFoundException, "repository not found", nil)
}

func (m *mockECRClient) SetRepositoryPolicyWithContext(ctx context.Context, input *ecr.SetRepositoryPolicyInput, opts ...request.Option) (*ecr.SetRepositoryPolicyOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, r := range tRepos {
		if aws.StringValue(input.RepositoryName) == aws.StringValue(r.RepositoryName) {
			return &ecr.SetRepositoryPolicyOutput{
				PolicyText:     input.PolicyText,
				RegistryId:     r.RegistryId,
				RepositoryName: r.RepositoryName,
			}, nil
		}
	}

	return nil, awserr.New(ecr.ErrCodeRepositoryNotFoundException, "repository not found", nil)
}

func TestECR_CreateRepository(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx   context.Context
		input *ecr.CreateRepositoryInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ecr.Repository
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:   context.TODO(),
				input: nil,
			},
			wantErr: true,
		},
		{
			name: "aws error",
			fields: fields{
				Service: newmockECRClient(t, awserr.New(ecr.ErrCodeEmptyUploadException, "bad request", nil)),
			},
			args: args{
				ctx:   context.TODO(),
				input: &ecr.CreateRepositoryInput{},
			},
			wantErr: true,
		},
		{
			name: "non-aws error",
			fields: fields{
				Service: newmockECRClient(t, errors.New("things blowing up!")),
			},
			args: args{
				ctx:   context.TODO(),
				input: &ecr.CreateRepositoryInput{},
			},
			wantErr: true,
		},
	}

	for _, repo := range tRepos {
		tests = append(tests, struct {
			name    string
			fields  fields
			args    args
			want    *ecr.Repository
			wantErr bool
		}{
			name: aws.StringValue(repo.RepositoryName),
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx: context.TODO(),
				input: &ecr.CreateRepositoryInput{
					EncryptionConfiguration:    repo.EncryptionConfiguration,
					ImageScanningConfiguration: repo.ImageScanningConfiguration,
					ImageTagMutability:         repo.ImageTagMutability,
					RepositoryName:             repo.RepositoryName,
				},
			},
			want: repo,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			got, err := e.CreateRepository(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECR.CreateRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ECR.CreateRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECR_ListRepositories(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "list repos",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{ctx: context.TODO()},
			want: []string{"carols/12DaysOfChristmas", "carols/SilentNight", "carols/FrostyTheSnowman", "carols/LittleDrummerBoy", "reindeer/rudolph", "reindeer/dasher", "reindeer/dancer"},
		},
		{
			name: "aws error",
			fields: fields{
				Service: newmockECRClient(t, awserr.New(ecr.ErrCodeEmptyUploadException, "bad request", nil)),
			},
			args:    args{ctx: context.TODO()},
			wantErr: true,
		},
		{
			name: "non-aws error",
			fields: fields{
				Service: newmockECRClient(t, errors.New("things blowing up!")),
			},
			args:    args{ctx: context.TODO()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			got, err := e.ListRepositories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECR.ListRepositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ECR.ListRepositories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECR_GetRepositories(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx      context.Context
		repoName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ecr.Repository
		wantErr bool
	}{
		{
			name: "empty input",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "",
			},
			wantErr: true,
		},
		{
			name: "unknown repository",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "somemissingrepo",
			},
			wantErr: true,
		},
		{
			name: "return many repositories",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "manyreposmatch",
			},
			wantErr: true,
		},
		{
			name: "aws error",
			fields: fields{
				Service: newmockECRClient(t, awserr.New(ecr.ErrCodeEmptyUploadException, "bad request", nil)),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "carols/JingleBells",
			},
			wantErr: true,
		},
		{
			name: "non-aws error",
			fields: fields{
				Service: newmockECRClient(t, errors.New("things blowing up!")),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "carols/JingleBells",
			},
			wantErr: true,
		},
	}

	for _, repo := range tRepos {
		tests = append(tests, struct {
			name    string
			fields  fields
			args    args
			want    *ecr.Repository
			wantErr bool
		}{
			name: aws.StringValue(repo.RepositoryName),
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: aws.StringValue(repo.RepositoryName),
			},
			want: repo,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			got, err := e.GetRepositories(tt.args.ctx, tt.args.repoName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECR.GetRepositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ECR.GetRepositories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECR_DeleteRepository(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx      context.Context
		repoName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ecr.Repository
		wantErr bool
	}{
		{
			name: "empty input",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "",
			},
			wantErr: true,
		},
		{
			name: "unknown repository",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "somemissingrepo",
			},
			wantErr: true,
		},
		{
			name: "aws error",
			fields: fields{
				Service: newmockECRClient(t, awserr.New(ecr.ErrCodeEmptyUploadException, "bad request", nil)),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "carols/JingleBells",
			},
			wantErr: true,
		},
		{
			name: "non-aws error",
			fields: fields{
				Service: newmockECRClient(t, errors.New("things blowing up!")),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: "carols/JingleBells",
			},
			wantErr: true,
		},
	}

	for _, repo := range tRepos {
		tests = append(tests, struct {
			name    string
			fields  fields
			args    args
			want    *ecr.Repository
			wantErr bool
		}{
			name: aws.StringValue(repo.RepositoryName),
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:      context.TODO(),
				repoName: aws.StringValue(repo.RepositoryName),
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			got, err := e.DeleteRepository(tt.args.ctx, tt.args.repoName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECR.DeleteRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ECR.DeleteRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECR_GetRepositoryTags(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx     context.Context
		repoArn string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*ecr.Tag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			got, err := e.GetRepositoryTags(tt.args.ctx, tt.args.repoArn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECR.GetRepositoryTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ECR.GetRepositoryTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECR_UpdateRepositoryTags(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx     context.Context
		repoArn string
		tags    []*ecr.Tag
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			if err := e.UpdateRepositoryTags(tt.args.ctx, tt.args.repoArn, tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("ECR.UpdateRepositoryTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestECR_SetImageScanningConfiguration(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx        context.Context
		repoName   string
		scanOnPush bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty input",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "",
				scanOnPush: true,
			},
			wantErr: true,
		},
		{
			name: "unknown repository",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "somemissingrepo",
				scanOnPush: true,
			},
			wantErr: true,
		},
		{
			name: "aws error",
			fields: fields{
				Service: newmockECRClient(t, awserr.New(ecr.ErrCodeEmptyUploadException, "bad request", nil)),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "carols/JingleBells",
				scanOnPush: true,
			},
			wantErr: true,
		},
		{
			name: "non-aws error",
			fields: fields{
				Service: newmockECRClient(t, errors.New("things blowing up!")),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "carols/JingleBells",
				scanOnPush: true,
			},
			wantErr: true,
		},
	}

	for _, repo := range tRepos {
		tests = append(tests, struct {
			name    string
			fields  fields
			args    args
			wantErr bool
		}{
			name: aws.StringValue(repo.RepositoryName),
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   aws.StringValue(repo.RepositoryName),
				scanOnPush: true,
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			if err := e.SetImageScanningConfiguration(tt.args.ctx, tt.args.repoName, tt.args.scanOnPush); (err != nil) != tt.wantErr {
				t.Errorf("ECR.SetImageScanningConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestECR_UpdateRepositoryPolicy(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx        context.Context
		repoName   string
		repoPolicy string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty repoName",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "",
				repoPolicy: "{}",
			},
			wantErr: true,
		},
		{
			name: "empty policy",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "carols/JingleBells",
				repoPolicy: "",
			},
			wantErr: true,
		},
		{
			name: "unknown repository",
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "somemissingrepo",
				repoPolicy: "{}",
			},
			wantErr: true,
		},
		{
			name: "aws error",
			fields: fields{
				Service: newmockECRClient(t, awserr.New(ecr.ErrCodeEmptyUploadException, "bad request", nil)),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "carols/JingleBells",
				repoPolicy: "{}",
			},
			wantErr: true,
		},
		{
			name: "non-aws error",
			fields: fields{
				Service: newmockECRClient(t, errors.New("things blowing up!")),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   "carols/JingleBells",
				repoPolicy: "{}",
			},
			wantErr: true,
		},
	}

	for _, repo := range tRepos {
		tests = append(tests, struct {
			name    string
			fields  fields
			args    args
			wantErr bool
		}{
			name: aws.StringValue(repo.RepositoryName),
			fields: fields{
				Service: newmockECRClient(t, nil),
			},
			args: args{
				ctx:        context.TODO(),
				repoName:   aws.StringValue(repo.RepositoryName),
				repoPolicy: "{}",
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			if err := e.UpdateRepositoryPolicy(tt.args.ctx, tt.args.repoName, tt.args.repoPolicy); (err != nil) != tt.wantErr {
				t.Errorf("ECR.UpdateRepositoryPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestECR_GetRepositoryPolicy(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx      context.Context
		repoName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECR{
				session:         tt.fields.session,
				Service:         tt.fields.Service,
				DefaultKMSKeyId: tt.fields.DefaultKMSKeyId,
			}
			got, err := e.GetRepositoryPolicy(tt.args.ctx, tt.args.repoName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECR.GetRepositoryPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ECR.GetRepositoryPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
