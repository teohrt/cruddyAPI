package service

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/entity"
)

type Service interface {
	CreateProfile(ctx context.Context, profile entity.Profile) (entity.CreateProfileResult, error)
	GetProfile(ctx context.Context, profileID string) (entity.Profile, error)
	UpdateProfile(ctx context.Context, profile entity.Profile) error
	DeleteProfile(ctx context.Context, profileID string) error
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

type ProfileNotFoundError struct {
	msg string
}

func (e ProfileNotFoundError) Error() string {
	return e.msg
}
