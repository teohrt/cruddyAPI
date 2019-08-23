package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/teohrt/cruddyAPI/dbclient/mock"
	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/entity"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	testCases := []struct {
		description           string
		profileID             string
		getItemOutputToReturn *dynamodb.GetItemOutput
		getItemErrorToReturn  error
		getItemReturnObject   interface{}
		expectedProfile       entity.Profile
		expectedErrorString   string
	}{
		{
			description:           "Happy path",
			profileID:             "123",
			getItemOutputToReturn: &dynamodb.GetItemOutput{},
			getItemErrorToReturn:  nil,
			getItemReturnObject:   entity.Profile{ID: "123"},
			expectedProfile:       entity.Profile{ID: "123"},
			expectedErrorString:   "",
		},
		{
			description:           "Profile not found",
			profileID:             "123",
			getItemOutputToReturn: &dynamodb.GetItemOutput{},
			getItemErrorToReturn:  nil,
			getItemReturnObject:   nil,
			expectedProfile:       entity.Profile{},
			expectedErrorString:   "Could not find profile associated with: 123",
		},
		{
			description:           "DB puked - QueryPK fails.",
			profileID:             "123",
			getItemOutputToReturn: nil,
			getItemErrorToReturn:  errors.New("puke"),
			getItemReturnObject:   nil,
			expectedProfile:       entity.Profile{},
			expectedErrorString:   "puke",
		},
		{
			description:           "Unmarshal error : Profile db had incompatible object",
			profileID:             "123",
			getItemOutputToReturn: &dynamodb.GetItemOutput{},
			getItemErrorToReturn:  nil,
			getItemReturnObject:   brokenProfileData{ID: false},
			expectedProfile:       entity.Profile{},
			expectedErrorString:   "UnmarshalTypeError: cannot unmarshal bool into Go value of type string",
		},
	}

	for _, tC := range testCases {
		logger := zerolog.New(os.Stdout)

		mockService := ServiceImpl{
			Client: dbclient.ClientImpl{
				DynamoDB: mock.DB{
					GetItemOutputToReturn: tC.getItemOutputToReturn,
					GetItemErrorToReturn:  tC.getItemErrorToReturn,
					GetItemReturnObject:   tC.getItemReturnObject,
				},
				Logger: &logger,
			},
			Logger: &logger,
		}

		profile, err := mockService.GetProfile(context.Background(), tC.profileID)
		if tC.expectedErrorString != "" {
			assert.Equal(t, tC.expectedErrorString, err.Error())
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tC.expectedProfile, profile)
	}
}

type brokenProfileData struct {
	ID bool `json: "id"`
}
