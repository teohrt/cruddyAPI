package dbclient

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

type Config struct {
	DynamoDBTableName  string `env:"DYNAMODB_TABLE_NAME"`
	AWSRegion          string `env:"AWS_SESSION_REGION"`
	AWSSessionEndpoint string `env:"AWS_SESSION_ENDPOINT"`
}

type Client interface {
	GetItem(ctx context.Context, valueName string, value string) (*map[string]*dynamodb.AttributeValue, error)
	UpsertItem(ctx context.Context, in interface{}) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, keyName string, value string) (*dynamodb.DeleteItemOutput, error)
}

type ClientImpl struct {
	DynamoDB  dynamodbiface.DynamoDBAPI
	TableName string
	Logger    *zerolog.Logger
}

func New(config *Config, logger *zerolog.Logger) Client {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(config.AWSRegion),
		Endpoint: aws.String(config.AWSSessionEndpoint),
	}))

	return ClientImpl{
		DynamoDB:  dynamodb.New(awsSession),
		TableName: config.DynamoDBTableName,
		Logger:    logger,
	}
}

func (db ClientImpl) GetItem(ctx context.Context, valueName string, value string) (*map[string]*dynamodb.AttributeValue, error) {
	result, err := db.DynamoDB.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			valueName: {
				S: aws.String(value),
			},
		},
	})

	if err != nil {
		db.Logger.Error().Err(err).Msg("GetItem failed")
		return nil, err
	}

	return &result.Item, nil
}

func (db ClientImpl) UpsertItem(ctx context.Context, in interface{}) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		db.Logger.Error().Err(err).Msg("Unable to MarshallMap for PutItem")
		return nil, err
	}

	return db.DynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.TableName),
	})
}

func (db ClientImpl) DeleteItem(ctx context.Context, keyName string, value string) (*dynamodb.DeleteItemOutput, error) {
	return db.DynamoDB.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(value),
			},
		},
	})
}
