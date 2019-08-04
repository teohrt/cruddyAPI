package testutils

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	MockTableName = "profiles"
)

type MockDB struct {
	dynamodbiface.DynamoDBAPI

	PutItemOutputToReturn *dynamodb.PutItemOutput
	PutItemErrorToReturn  error

	GetItemOutputToReturn *dynamodb.GetItemOutput
	GetItemReturnObject   interface{}
	GetItemErrorToReturn  error

	DeleteItemOutputToReturn *dynamodb.DeleteItemOutput
	DeleteItemErrorToReturn  error
}

func (m MockDB) PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput, options ...request.Option) (*dynamodb.PutItemOutput, error) {
	return m.PutItemOutputToReturn, m.PutItemErrorToReturn
}

func (m MockDB) GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput, options ...request.Option) (*dynamodb.GetItemOutput, error) {
	o := m.GetItemOutputToReturn
	if o != nil {
		o.Item = map[string]*dynamodb.AttributeValue{}
	}

	if m.GetItemReturnObject != nil {
		profile, _ := dynamodbattribute.MarshalMap(m.GetItemReturnObject)
		o.Item = profile
	}
	return o, m.GetItemErrorToReturn
}

func (m MockDB) DeleteItemWithContext(ctx context.Context, input *dynamodb.DeleteItemInput, options ...request.Option) (*dynamodb.DeleteItemOutput, error) {
	return m.DeleteItemOutputToReturn, m.DeleteItemErrorToReturn
}
