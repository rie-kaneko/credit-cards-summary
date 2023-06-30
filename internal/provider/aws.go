package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Service struct {
	Session      *session.Session
	SesClient    *ses.SES
	DynamoClient *dynamodb.DynamoDB
}

func NewService(region string) (*Service, error) {
	sess, err := createSession(region)
	if err != nil {
		return nil, err
	}

	return &Service{
		Session:      sess,
		SesClient:    createSESClient(sess),
		DynamoClient: createDynamoClient(sess),
	}, nil
}

func createSession(region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
}

func createSESClient(sess *session.Session) *ses.SES {
	return ses.New(sess)
}

func createDynamoClient(sess *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(sess)
}
