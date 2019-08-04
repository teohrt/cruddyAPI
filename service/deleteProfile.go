package service

import (
	"context"
)

func (svc serviceImpl) DeleteProfile(ctx context.Context, profileID string) error {
	_, err := svc.Client.DeleteItem(ctx, "id", profileID)
	if err != nil {
		svc.Logger.Error().Err(err).Msg("DeleteItem failed")
		return err
	}

	return nil
}
