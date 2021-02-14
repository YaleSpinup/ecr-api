package api

import (
	"github.com/YaleSpinup/ecr-api/ecr"
	"github.com/YaleSpinup/ecr-api/iam"
)

type ecrOrchestrator struct {
	client ecr.ECR
	org    string
}

func newEcrOrchestrator(client ecr.ECR, org string) *ecrOrchestrator {
	return &ecrOrchestrator{
		client: client,
		org:    org,
	}
}

type iamOrchestrator struct {
	client iam.IAM
	org    string
}

func newIamOrchestrator(client iam.IAM, org string) *iamOrchestrator {
	return &iamOrchestrator{
		client: client,
		org:    org,
	}
}
