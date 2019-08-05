package service

import (
	"context"

	"github.com/teohrt/cruddyAPI/entity"
)

func (svc ServiceImpl) UpdateProfile(ctx context.Context, profile entity.Profile) error {
	_, err := svc.GetProfile(ctx, profile.ID)
	if err != nil {
		svc.Logger.Error().Err(err).Msg("UpdateProfile: GetProfile failed")
		return err
	}

	_, err = svc.Client.UpsertItem(ctx, profile)
	if err != nil {
		svc.Logger.Printf("%v", profile.ID)
		svc.Logger.Error().Err(err).Msg("UpdateProfile: UpsertItem failed")
		return err
	}

	return nil
}
