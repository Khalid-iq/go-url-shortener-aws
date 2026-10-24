package main

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoStore keeps links in a DynamoDB table, so they survive restarts and
// are shared by every running copy of the app.
type DynamoStore struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoStore(client *dynamodb.Client, table string) *DynamoStore {
	return &DynamoStore{client: client, table: table}
}

// Each call gets a deadline so a slow DynamoDB call can't hang a request.
func requestCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func (s *DynamoStore) Save(code, url string) error {
	ctx, cancel := requestCtx()
	defer cancel()

	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.table),
		Item: map[string]types.AttributeValue{
			"code": &types.AttributeValueMemberS{Value: code},
			"url":  &types.AttributeValueMemberS{Value: url},
		},
		// Never overwrite an existing code, even if two requests race.
		ConditionExpression: aws.String("attribute_not_exists(code)"),
	})
	return err
}

func (s *DynamoStore) Get(code string) (string, error) {
	ctx, cancel := requestCtx()
	defer cancel()

	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.table),
		Key: map[string]types.AttributeValue{
			"code": &types.AttributeValueMemberS{Value: code},
		},
	})
	if err != nil {
		return "", err
	}
	if out.Item == nil {
		return "", ErrNotFound
	}
	v, ok := out.Item["url"].(*types.AttributeValueMemberS)
	if !ok {
		return "", errors.New("stored item has no url")
	}
	return v.Value, nil
}

func (s *DynamoStore) Exists(code string) bool {
	_, err := s.Get(code)
	return err == nil
}
