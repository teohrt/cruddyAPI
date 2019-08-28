package service

import (
	"context"

	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/utils"
)

func (svc ServiceImpl) UpdateProfile(ctx context.Context, profileData entity.ProfileData, profileID string) error {
	ctx, seg := utils.StartXraySegment(ctx, "UpdateProfile service")
	emailHash := utils.Hash(profileData.Email)
	if profileID != emailHash {
		err := EmailIncsonsistentWithProfileIDError{"Email inconsistent with ProfileID"}
		seg.Close(err)
		return err
	}

	_, err := svc.GetProfile(ctx, profileID)
	if err != nil {
		seg.Close(err)
		svc.Logger.Error().Err(err).Msg("UpdateProfile: GetProfile failed")
		return err
	}

	updateProfile := entity.Profile{
		ID:          profileID,
		ProfileData: profileData,
	}

	_, err = svc.Client.UpsertItem(ctx, updateProfile)
	if err != nil {
		seg.Close(err)
		svc.Logger.Printf("%v", profileID)
		svc.Logger.Error().Err(err).Msg("UpdateProfile: UpsertItem failed")
		return err
	}

	seg.Close(nil)
	return nil
}
