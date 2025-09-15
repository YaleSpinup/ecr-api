package ecr

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

// TestGetImageScanFindingsCompatibility tests the GetImageScanFindings function
// with various scenarios to ensure compatibility with AWS SDK v1.49.16
func TestGetImageScanFindingsCompatibility(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		tag     string
		mockErr error
		wantErr bool
		errMsg  string
	}{
		{
			name:    "successful scan findings retrieval",
			repo:    "spinup-000941/spinup-002e8c-wsic-repo",
			tag:     "latest",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "scan not found - image never scanned",
			repo:    "spinup-000941/spinup-002e8c-wsic-repo",
			tag:     "latest",
			mockErr: awserr.New("ScanNotFoundException", "Image scan does not exist for the image", nil),
			wantErr: true,
			errMsg:  "ScanNotFoundException",
		},
		{
			name:    "image not found",
			repo:    "non-existent-repo",
			tag:     "latest",
			mockErr: awserr.New("ImageNotFoundException", "The image requested does not exist", nil),
			wantErr: true,
			errMsg:  "ImageNotFoundException",
		},
		{
			name:    "repository not found",
			repo:    "non-existent-repo",
			tag:     "latest",
			mockErr: awserr.New("RepositoryNotFoundException", "The repository does not exist", nil),
			wantErr: true,
			errMsg:  "RepositoryNotFoundException",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockECRScanClient{
				scanErr: tt.mockErr,
			}

			e := &ECR{
				Service: mockClient,
			}

			findings, err := e.GetImageScanFindings(context.Background(), tt.repo, tt.tag)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if findings == nil {
					t.Errorf("expected findings but got nil")
				}
			}
		})
	}
}

// TestGetImageScanFindingsByImageDigestCompatibility tests the GetImageScanFindingsByImageDigest function
func TestGetImageScanFindingsByImageDigestCompatibility(t *testing.T) {
	tests := []struct {
		name        string
		repo        string
		imageDigest string
		mockErr     error
		wantErr     bool
	}{
		{
			name:        "successful scan findings by digest",
			repo:        "spinup-000941/spinup-002e8c-wsic-repo",
			imageDigest: "sha256:1234567890abcdef",
			mockErr:     nil,
			wantErr:     false,
		},
		{
			name:        "scan not found by digest",
			repo:        "spinup-000941/spinup-002e8c-wsic-repo",
			imageDigest: "sha256:1234567890abcdef",
			mockErr:     awserr.New("ScanNotFoundException", "Image scan does not exist", nil),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockECRScanClient{
				scanErr: tt.mockErr,
			}

			e := &ECR{
				Service: mockClient,
			}

			output, err := e.GetImageScanFindingsByImageDigest(context.Background(), tt.repo, tt.imageDigest)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if output == nil {
					t.Errorf("expected output but got nil")
				}
			}
		})
	}
}

// mockECRScanClient is a mock implementation of the ECR client for scan tests
type mockECRScanClient struct {
	ecriface.ECRAPI
	scanErr error
}

func (m *mockECRScanClient) DescribeImageScanFindingsWithContext(ctx context.Context, input *ecr.DescribeImageScanFindingsInput, opts ...request.Option) (*ecr.DescribeImageScanFindingsOutput, error) {
	if m.scanErr != nil {
		return nil, m.scanErr
	}

	return &ecr.DescribeImageScanFindingsOutput{
		ImageId: input.ImageId,
		ImageScanFindings: &ecr.ImageScanFindings{
			FindingSeverityCounts: map[string]*int64{
				"HIGH":         aws.Int64(2),
				"MEDIUM":       aws.Int64(5),
				"LOW":          aws.Int64(10),
				"INFORMATIONAL": aws.Int64(3),
			},
			Findings: []*ecr.ImageScanFinding{},
		},
		ImageScanStatus: &ecr.ImageScanStatus{
			Status:      aws.String("COMPLETE"),
			Description: aws.String("Scan complete"),
		},
		RepositoryName: input.RepositoryName,
	}, nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}