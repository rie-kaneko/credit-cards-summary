package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
	"rie-kaneko/credit-cards-summary/config"
)

type Service struct {
	Session      *session.Session
	SesClient    *ses.SES
	DynamoClient *dynamodb.DynamoDB
}

func NewService(conf config.AWSConfig) (*Service, error) {
	sess, err := createSession(conf)
	if err != nil {
		return nil, err
	}

	return &Service{
		Session:      sess,
		SesClient:    createSESClient(sess),
		DynamoClient: createDynamoClient(sess),
	}, nil
}

func createSession(conf config.AWSConfig) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: credentials.NewStaticCredentials(conf.AccessKey, conf.SecretKey, ""),
	})
}

func createSESClient(sess *session.Session) *ses.SES {
	return ses.New(sess)
}

func createDynamoClient(sess *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(sess)
}
