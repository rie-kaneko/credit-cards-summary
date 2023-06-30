package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"time"
)

const (
	transactionTableName = "transactions"
)

type Status string

const (
	Succeeded Status = "succeeded"
	Failed    Status = "failed"
)

type Transaction struct {
	ID         string    `json:"id"`
	Sender     string    `json:"sender"`
	ReceiverID string    `json:"receiver_id"`
	Receiver   string    `json:"receiver"`
	Time       time.Time `json:"time"`
	Status     Status    `json:"status"`
}

func (s *Service) PostTransaction(u Transaction) error {
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(transactionTableName),
	}

	_, err = s.DynamoClient.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}
