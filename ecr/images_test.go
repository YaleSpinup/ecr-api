package ecr

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

func TestECR_ListImages(t *testing.T) {
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
		want    []*ecr.ImageIdentifier
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
			got, err := e.ListImages(tt.args.ctx, tt.args.repoName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECR.ListImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ECR.ListImages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECR_GetImages(t *testing.T) {
	type fields struct {
		session         *session.Session
		Service         ecriface.ECRAPI
		DefaultKMSKeyId string
	}
	type args struct {
		ctx      context.Context
		repoName string
		imageIds []*ecr.ImageIdentifier
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*ecr.ImageDetail
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
			got, err := e.GetImages(tt.args.ctx, tt.args.repoName, tt.args.imageIds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ECR.GetImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ECR.GetImages() = %v, want %v", got, tt.want)
			}
		})
	}
}
