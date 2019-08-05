package service

import (
	"context"

	"github.com/teohrt/cruddyAPI/entity"
)

func (svc ServiceImpl) CreateProfile(ctx context.Context, profile entity.Profile) (entity.CreateProfileResult, error) {
	_, err := svc.GetProfile(ctx, profile.ID)
	if err != nil {
		switch err.(type) {
		case ProfileNotFoundError: // Profile has not been previous created for this ID. We can proceed with creation
			_, err = svc.Client.UpsertItem(ctx, profile)
			if err != nil {
				svc.Logger.Printf("%v", profile.ID)
				svc.Logger.Error().Err(err).Msg("Upsert profile failed")
				return entity.CreateProfileResult{}, err
			}
			return entity.CreateProfileResult{ProfileID: profile.ID}, nil

		default:
			svc.Logger.Error().Err(err).Msg("GetProfile: failed retrieving profile from ID: " + profile.ID)
			return entity.CreateProfileResult{}, err
		}
	}

	return entity.CreateProfileResult{}, ProfileAlreadyExistsError{"Can not create profile. Already exists"}
}
