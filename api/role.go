package api

import (
	"context"
	"fmt"

	"github.com/YaleSpinup/ecr-api/session"
	stsSvc "github.com/YaleSpinup/ecr-api/sts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (s *server) assumeRole(ctx context.Context, id, role, inlinePolicy string, policyArns ...string) (*session.Session, error) {
	log.Infof("assuming role %s", role)

	stsService := stsSvc.New(stsSvc.WithSession(s.session.Session))

	name := fmt.Sprintf("spinup-%s-ecr-api-%s", s.org, uuid.New())

	input := sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(900),
		RoleArn:         aws.String(role),
		RoleSessionName: aws.String(name),
		Tags: []*sts.Tag{
			{
				Key:   aws.String("spinup:org"),
				Value: aws.String(s.org),
			},
		},
	}

	if id != "" {
		input.SetExternalId(id)
	}

	if inlinePolicy != "" {
		input.SetPolicy(inlinePolicy)
	}

	log.Debugf("assuming role %s with input %+v", role, input)

	out, err := stsService.AssumeRole(ctx, &input)
	if err != nil {
		log.Errorf("got: %s", err)
		return nil, err
	}

	akid := aws.StringValue(out.Credentials.AccessKeyId)

	log.Infof("got temporary creds %s, expiration: %s", akid, aws.TimeValue(out.Credentials.Expiration).String())

	sess := session.New(
		session.WithCredentials(
			akid,
			aws.StringValue(out.Credentials.SecretAccessKey),
			aws.StringValue(out.Credentials.SessionToken),
		),
		session.WithRegion("us-east-1"),
	)

	return &sess, nil
}
