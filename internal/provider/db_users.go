package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	userTableName = "users"
	userKeyName   = "id"
)

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	EmailAddress string `json:"email_address"`
}

func (s *Service) GetUser(id string) (*User, error) {
	result, err := s.DynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(userTableName),
		Key: map[string]*dynamodb.AttributeValue{
			userKeyName: {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return convertToUser(result)
}

func (s *Service) PostUser(u User) error {
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(userTableName),
	}

	_, err = s.DynamoClient.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func convertToUser(aUser *dynamodb.GetItemOutput) (*User, error) {
	var user User
	err := dynamodbattribute.UnmarshalMap(aUser.Item, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
