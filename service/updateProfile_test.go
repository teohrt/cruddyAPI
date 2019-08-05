package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/testutils"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProfile(t *testing.T) {
	testCases := []struct {
		description           string
		profile               entity.Profile
		GetItemOutputToReturn *dynamodb.GetItemOutput
		GetItemErrorToReturn  error
		GetItemReturnObject   interface{}
		PutItemOutputToReturn *dynamodb.PutItemOutput
		PutItemErrorToReturn  error
		expectedErrorString   string
	}{
		{
			description:           "Happy path",
			profile:               entity.Profile{ID: "123"},
			GetItemOutputToReturn: &dynamodb.GetItemOutput{},
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   entity.Profile{ID: "123"},
			expectedErrorString:   "",
		},
		{
			description:           "Unmarshal error : Profile db had incompatible object",
			profile:               entity.Profile{ID: "123"},
			GetItemOutputToReturn: &dynamodb.GetItemOutput{},
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   brokenProfileData{ID: false},
			expectedErrorString:   "UnmarshalTypeError: cannot unmarshal bool into Go value of type string",
		},
		{
			description:           "ID not found",
			profile:               entity.Profile{ID: "123"},
			GetItemOutputToReturn: &dynamodb.GetItemOutput{},
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   entity.Profile{ID: ""},
			expectedErrorString:   "Could not find profile associated with: 123",
		},
		{
			description:           "Profile UpsertItem puked",
			profile:               entity.Profile{ID: "123"},
			GetItemOutputToReturn: &dynamodb.GetItemOutput{},
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   entity.Profile{ID: "123"},
			PutItemOutputToReturn: nil,
			PutItemErrorToReturn:  errors.New("puke"),
			expectedErrorString:   "puke",
		},
	}

	for _, tC := range testCases {
		logger := zerolog.New(os.Stdout)

		mockService := ServiceImpl{
			Client: dbclient.ClientImpl{
				DynamoDB: testutils.MockDB{
					GetItemOutputToReturn: tC.GetItemOutputToReturn,
					GetItemErrorToReturn:  tC.GetItemErrorToReturn,
					GetItemReturnObject:   tC.GetItemReturnObject,
					PutItemOutputToReturn: tC.PutItemOutputToReturn,
					PutItemErrorToReturn:  tC.PutItemErrorToReturn,
				},
				Logger: &logger,
			},
			Logger: &logger,
		}

		err := mockService.UpdateProfile(context.Background(), tC.profile)

		if tC.expectedErrorString != "" {
			assert.Error(t, err)
			assert.Equal(t, tC.expectedErrorString, err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
