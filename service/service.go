package service

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/entity"
)

type Service interface {
	GetProfileService(ctx context.Context, profileID int) (*entity.Profile, error)
	CreateProfileService(ctx context.Context, profile *entity.Profile) (*entity.CreateProfileResult, error)
}

type serviceImpl struct {
	Client    dbclient.Client
	Logger    *zerolog.Logger
	TableName string
}

func New(config *dbclient.DBConfig) Service {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	client := dbclient.New(config, &logger)
	return serviceImpl{
		Client:    client,
		Logger:    &logger,
		TableName: config.DynamoDBTableName,
	}
}
