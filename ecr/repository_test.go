package ecr

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"
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

func TestCreateRepository(t *testing.T) {
	r := ECR{Service: newmockECRClient(t, nil)}

	// nil input
	if _, err := r.CreateRepository(context.TODO(), nil); err == nil {
		t.Error("expected error, got nil")
	}

	for _, repo := range tRepos {
		t.Logf("testing create repository %s", aws.StringValue(repo.RepositoryName))

		out, err := r.CreateRepository(context.TODO(), &ecr.CreateRepositoryInput{
			EncryptionConfiguration:    repo.EncryptionConfiguration,
			ImageScanningConfiguration: repo.ImageScanningConfiguration,
			ImageTagMutability:         repo.ImageTagMutability,
			RepositoryName:             repo.RepositoryName,
		})

		if err != nil {
			t.Errorf("expected nil error, got %s", err)
		}

		if !awsutil.DeepEqual(out, repo) {
			t.Errorf("expected %s, got %s", awsutil.Prettify(repo), awsutil.Prettify(out))
		}
	}

	r.Service.(*mockECRClient).err = awserr.New(ecr.ErrCodeEmptyUploadException, "bad request", nil)
	_, err := r.CreateRepository(context.TODO(), &ecr.CreateRepositoryInput{})
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrBadRequest {
			t.Errorf("expected error code %s, got: %s", apierror.ErrBadRequest, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}

	r.Service.(*mockECRClient).err = awserr.New(ecr.ErrCodeServerException, "internal error", nil)
	_, err = r.CreateRepository(context.TODO(), &ecr.CreateRepositoryInput{})
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrInternalError {
			t.Errorf("expected error code %s, got: %s", apierror.ErrInternalError, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}

	r.Service.(*mockECRClient).err = awserr.New(ecr.ErrCodeRepositoryNotFoundException, "not found", nil)
	_, err = r.CreateRepository(context.TODO(), &ecr.CreateRepositoryInput{})
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrNotFound {
			t.Errorf("expected error code %s, got: %s", apierror.ErrNotFound, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}

	r.Service.(*mockECRClient).err = awserr.New(ecr.ErrCodeRepositoryAlreadyExistsException, "in use", nil)
	_, err = r.CreateRepository(context.TODO(), &ecr.CreateRepositoryInput{})
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrConflict {
			t.Errorf("expected error code %s, got: %s", apierror.ErrConflict, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}

	// test non-aws error
	r.Service.(*mockECRClient).err = errors.New("things blowing up!")
	_, err = r.CreateRepository(context.TODO(), &ecr.CreateRepositoryInput{})
	if aerr, ok := err.(apierror.Error); ok {
		if aerr.Code != apierror.ErrInternalError {
			t.Errorf("expected error code %s, got: %s", apierror.ErrInternalError, aerr.Code)
		}
	} else {
		t.Errorf("expected apierror.Error, got: %s", reflect.TypeOf(err).String())
	}
}

func TestListRepositories(t *testing.T) {
	t.Log("todo")
}

func TestDeleteRepository(t *testing.T) {
	t.Log("todo")
}

func TestGetRepositoryTags(t *testing.T) {
	t.Log("todo")
}

func TestUpdateRepositoryTags(t *testing.T) {
	t.Log("todo")
}

func TestSetImageScanningConfiguration(t *testing.T) {
	t.Log("todo")
}

func TestGetRepositories(t *testing.T) {
	t.Log("todo")
}
