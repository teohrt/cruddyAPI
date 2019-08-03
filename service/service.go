package service

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/entity"
)

type Service interface {
	CreateProfile(ctx context.Context, profile entity.Profile) (*entity.CreateProfileResult, error)
	GetProfile(ctx context.Context, profileID int) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, profile entity.Profile) error
	DeleteProfile(ctx context.Context, profileID int) error
}

type serviceImpl struct {
	Client dbclient.Client
	Logger *zerolog.Logger
}

func New(config *dbclient.DBConfig) Service {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	client := dbclient.New(config, &logger)
	return serviceImpl{
		Client: client,
		Logger: &logger,
	}
}
