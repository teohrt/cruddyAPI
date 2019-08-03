package service

import (
	"context"

	"github.com/teohrt/cruddyAPI/entity"
)

func (svc serviceImpl) CreateProfile(ctx context.Context, profile entity.Profile) (*entity.CreateProfileResult, error) {
	_, err := svc.Client.UpsertItem(ctx, profile)
	if err != nil {
		svc.Logger.Printf("%v", profile.ID)
		svc.Logger.Error().Err(err).Msg("Adding profile failed")
		return &entity.CreateProfileResult{}, err
	}

	return &entity.CreateProfileResult{profile.ID}, nil
}
