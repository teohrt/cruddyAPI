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

const (
	EMAIL      = "foo@bar.com"
	EMAIL_HASH = "0c7e6a405862e402eb76a70f8a26fc732d07c32931e9fae9ab1582911d2e8a3b"
)

func TestUpdateProfile(t *testing.T) {
	testCases := []struct {
		description           string
		profileID             string
		profileData           entity.ProfileData
		GetItemOutputToReturn *dynamodb.GetItemOutput
		GetItemErrorToReturn  error
		GetItemReturnObject   interface{}
		PutItemOutputToReturn *dynamodb.PutItemOutput
		PutItemErrorToReturn  error
		expectedErrorString   string
	}{
		{
			description:           "Happy path",
			profileID:             EMAIL_HASH,
			profileData:           entity.ProfileData{Email: EMAIL},
			GetItemOutputToReturn: &dynamodb.GetItemOutput{},
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   entity.Profile{ID: EMAIL_HASH},
			expectedErrorString:   "",
		},
		{
			description:           "Unmarshal error : Profile db had incompatible object",
			profileID:             EMAIL_HASH,
			profileData:           entity.ProfileData{Email: "foo@bar.com"},
			GetItemOutputToReturn: &dynamodb.GetItemOutput{},
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   brokenProfileData{ID: false},
			expectedErrorString:   "UnmarshalTypeError: cannot unmarshal bool into Go value of type string",
		},
		{
			description:           "ID not found",
			profileID:             EMAIL_HASH,
			profileData:           entity.ProfileData{Email: EMAIL},
			GetItemOutputToReturn: &dynamodb.GetItemOutput{},
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   entity.Profile{ID: ""},
			expectedErrorString:   "Could not find profile associated with: " + EMAIL_HASH,
		},
		{
			description:           "Profile UpsertItem puked",
			profileID:             EMAIL_HASH,
			profileData:           entity.ProfileData{Email: EMAIL},
			GetItemOutputToReturn: &dynamodb.GetItemOutput{},
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   entity.Profile{ID: EMAIL_HASH},
			PutItemOutputToReturn: nil,
			PutItemErrorToReturn:  errors.New("puke"),
			expectedErrorString:   "puke",
		},
		{
			description:           "UpdateProfile puked - Email inconsistent with pid",
			profileID:             "bad-profileID",
			profileData:           entity.ProfileData{Email: EMAIL},
			GetItemOutputToReturn: nil,
			GetItemErrorToReturn:  nil,
			GetItemReturnObject:   nil,
			PutItemOutputToReturn: nil,
			PutItemErrorToReturn:  nil,
			expectedErrorString:   "Email inconsistent with ProfileID",
		},
	}

	for _, tC := range testCases {
		logger := zerolog.New(os.Stdout)

		mockService := ServiceImpl{
			Client: dbclient.ClientImpl{
				DynamoDB: mock.DB{
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

		err := mockService.UpdateProfile(context.Background(), tC.profileData, tC.profileID)

		if tC.expectedErrorString != "" {
			assert.Error(t, err)
			assert.Equal(t, tC.expectedErrorString, err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
