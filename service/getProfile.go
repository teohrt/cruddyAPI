package service

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/teohrt/cruddyAPI/entity"
)

// TODO
func (svc serviceImpl) GetProfile(ctx context.Context, profileID string) (entity.Profile, error) {
	item, err := svc.Client.GetItem(ctx, "ProfileID", profileID)
	if err != nil {
		svc.Logger.Error().Err(err).Msg("Searching for ProfileID failed")
		return entity.Profile{}, err
	}

	profile := new(entity.Profile)
	err = dynamodbattribute.UnmarshalMap(*item, &profile)
	if err != nil {
		svc.Logger.Warn().Err(err).Msg("Unmarshalling Map failed. Item from db doesn't match dto. PID: " + profileID)
		return entity.Profile{}, err
	}

	if profile.ID == "" {
		svc.Logger.Debug().Err(err).Msg("Could not find profile.")
		return entity.Profile{}, ProfileNotFoundError{"Could not find profile associated with: " + profileID}
	}

	return *profile, nil
}
