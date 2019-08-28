package service

import (
	"context"

	"github.com/teohrt/cruddyAPI/utils"
)

func (svc ServiceImpl) DeleteProfile(ctx context.Context, profileID string) error {
	ctx, seg := utils.StartXraySegment(ctx, "DeleteProfile service")
	if _, err := svc.GetProfile(ctx, profileID); err != nil {
		seg.Close(err)
		svc.Logger.Error().Err(err).Msg("GetProfile failed")
		return err
	}

	_, err := svc.Client.DeleteItem(ctx, "id", profileID)
	if err != nil {
		seg.Close(err)
		svc.Logger.Error().Err(err).Msg("DeleteItem failed")
		return err
	}

	seg.Close(nil)
	return nil
}
