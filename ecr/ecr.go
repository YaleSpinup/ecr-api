package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
)

type ECR struct {
	session         *session.Session
	Service         ecriface.ECRAPI
	DefaultKMSKeyId string
}

type ECROption func(*ECR)

func New(opts ...ECROption) ECR {
	e := ECR{}

	for _, opt := range opts {
		opt(&e)
	}

	if e.session != nil {
		e.Service = ecr.New(e.session)
	}

	return e
}

func WithSession(sess *session.Session) ECROption {
	return func(e *ECR) {
		log.Debug("using aws session")
		e.session = sess
	}
}

func WithCredentials(key, secret, token, region string) ECROption {
	return func(e *ECR) {
		log.Debugf("creating new session with key id %s in region %s", key, region)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, token),
			Region:      aws.String(region),
		}))
		e.session = sess
	}
}

func WithDefaultKMSKeyId(keyId string) ECROption {
	return func(e *ECR) {
		log.Debugf("using default kms keyid %s", keyId)
		e.DefaultKMSKeyId = keyId
	}
}
