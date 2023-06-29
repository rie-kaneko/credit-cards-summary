package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Service struct {
	Session   *session.Session
	SesClient *ses.SES
}

func NewService() (*Service, error) {
	sess, err := createSession()
	if err != nil {
		return nil, err
	}

	return &Service{
		Session:   sess,
		SesClient: createSESClient(sess),
	}, nil
}

func createSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "my-name",
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func createSESClient(sess *session.Session) *ses.SES {
	return ses.New(sess)
}
