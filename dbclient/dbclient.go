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

type clientImpl struct {
	dynamoDB  dynamodbiface.DynamoDBAPI
	tableName string
	logger    *zerolog.Logger
}

func New(config *Config, logger *zerolog.Logger) Client {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(config.AWSRegion),
		Endpoint: aws.String(config.AWSSessionEndpoint),
	}))

	return clientImpl{
		dynamoDB:  dynamodb.New(awsSession),
		tableName: config.DynamoDBTableName,
		logger:    logger,
	}
}

func (db clientImpl) GetItem(ctx context.Context, valueName string, value string) (*map[string]*dynamodb.AttributeValue, error) {
	result, err := db.dynamoDB.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(db.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			valueName: {
				S: aws.String(value),
			},
		},
	})

	if err != nil {
		db.logger.Error().Err(err).Msg("GetItem failed")
		return nil, err
	}

	return &result.Item, nil
}

func (db clientImpl) UpsertItem(ctx context.Context, in interface{}) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		db.logger.Error().Err(err).Msg("Unable to MarshallMap for PutItem")
		return nil, err
	}

	return db.dynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.tableName),
	})
}

func (db clientImpl) DeleteItem(ctx context.Context, keyName string, value string) (*dynamodb.DeleteItemOutput, error) {
	return db.dynamoDB.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(db.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(value),
			},
		},
	})
}
