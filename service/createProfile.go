package service

import (
	"context"

	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/utils"
)

func (svc ServiceImpl) CreateProfile(ctx context.Context, profileData entity.ProfileData) (entity.CreateProfileResult, error) {
	ctx, seg := utils.StartXraySegment(ctx, "CreateProfile service")
	profileID := utils.Hash(profileData.Email)

	_, err := svc.GetProfile(ctx, profileID)
	if err != nil {
		seg.Close(err)
		switch err.(type) {
		case ProfileNotFoundError: // Profile has not been previously created for this ID. We can proceed with creation
			_, err = svc.Client.UpsertItem(ctx, entity.Profile{
				ID:          profileID,
				ProfileData: profileData,
			})
			if err != nil {
				svc.Logger.Printf("%v", profileID)
				svc.Logger.Error().Err(err).Msg("Upsert profile failed")
				return entity.CreateProfileResult{}, err
			}
			return entity.CreateProfileResult{ProfileID: profileID}, nil

		default:
			svc.Logger.Error().Err(err).Msg("GetProfile: failed retrieving profile from ID: " + profileID)
			return entity.CreateProfileResult{}, err
		}
	}

	seg.Close(nil)
	return entity.CreateProfileResult{}, ProfileAlreadyExistsError{"Can not create profile. Already exists"}
}
