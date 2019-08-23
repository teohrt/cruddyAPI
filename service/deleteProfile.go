package service

import (
	"context"
)

func (svc ServiceImpl) DeleteProfile(ctx context.Context, profileID string) error {
	if _, err := svc.GetProfile(ctx, profileID); err != nil {
		svc.Logger.Error().Err(err).Msg("GetProfile failed")
		return err
	}

	_, err := svc.Client.DeleteItem(ctx, "id", profileID)
	if err != nil {
		svc.Logger.Error().Err(err).Msg("DeleteItem failed")
		return err
	}

	return nil
}
