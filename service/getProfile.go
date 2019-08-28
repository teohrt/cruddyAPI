package service

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/utils"
)

func (svc ServiceImpl) GetProfile(ctx context.Context, profileID string) (entity.Profile, error) {
	ctx, seg := utils.StartXraySegment(ctx, "GetProfile service")
	item, err := svc.Client.GetItem(ctx, "id", profileID)
	if err != nil {
		seg.Close(err)
		svc.Logger.Error().Err(err).Msg("Searching for ProfileID failed")
		return entity.Profile{}, err
	}

	profile := new(entity.Profile)
	err = dynamodbattribute.UnmarshalMap(*item, &profile)
	if err != nil {
		seg.Close(err)
		svc.Logger.Warn().Err(err).Msg("Unmarshalling Map failed. Item from db doesn't match dto. PID: " + profileID)
		return entity.Profile{}, err
	}

	if profile.ID == "" {
		err = ProfileNotFoundError{"Could not find profile associated with: " + profileID}
		seg.Close(err)
		svc.Logger.Debug().Err(err).Msg("Could not find profile.")
		return entity.Profile{}, err
	}

	seg.Close(nil)
	return *profile, nil
}
