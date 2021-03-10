package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	log "github.com/sirupsen/logrus"
)

type IAM struct {
	session *session.Session
	Service iamiface.IAMAPI
}

type IAMOption func(*IAM)

func New(opts ...IAMOption) IAM {
	i := IAM{}

	for _, opt := range opts {
		opt(&i)
	}

	if i.session != nil {
		i.Service = iam.New(i.session)
	}

	return i
}

func WithSession(sess *session.Session) IAMOption {
	return func(i *IAM) {
		log.Debug("using aws session")
		i.session = sess
	}
}

func WithCredentials(key, secret, token, region string) IAMOption {
	return func(i *IAM) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		i.session = sess
	}
}
