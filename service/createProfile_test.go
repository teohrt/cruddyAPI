package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/dbclient/mock"
	"github.com/teohrt/cruddyAPI/entity"
)

func TestCreateProfile(t *testing.T) {
	testCases := []struct {
		description           string
		profile               entity.ProfileData
		putItemOutputToReturn *dynamodb.PutItemOutput
		putItemErrorToReturn  error
		getItemOutput         *dynamodb.GetItemOutput
		getItemReturnObject   interface{}
		getItemErrorToReturn  error
		expectedErrorString   string
		expectedResultString  string
	}{
		{
			description:           "Happy Path",
			profile:               entity.ProfileData{Email: EMAIL},
			putItemOutputToReturn: &dynamodb.PutItemOutput{},
			putItemErrorToReturn:  nil,
			getItemOutput:         nil,
			getItemReturnObject:   nil,
			getItemErrorToReturn:  ProfileNotFoundError{"profile not found"},
			expectedErrorString:   "",
			expectedResultString:  fmt.Sprintf("{\"ProfileID\":\"%s\"}\n", EMAIL_HASH),
		},
		{
			description:           "DB error - UpsertItem Puked",
			profile:               entity.ProfileData{Email: EMAIL},
			putItemOutputToReturn: &dynamodb.PutItemOutput{},
			putItemErrorToReturn:  errors.New("puke"),
			getItemOutput:         nil,
			getItemReturnObject:   nil,
			getItemErrorToReturn:  ProfileNotFoundError{"profile not found"},
			expectedErrorString:   "puke",
			expectedResultString:  fmt.Sprintf("{\"ProfileID\":\"%s\"}\n", ""),
		},
		{
			description:           "DB error - GetItem Puked",
			profile:               entity.ProfileData{Email: "foo@bar.com"},
			putItemOutputToReturn: nil,
			putItemErrorToReturn:  nil,
			getItemOutput:         &dynamodb.GetItemOutput{},
			getItemReturnObject:   entity.Profile{ID: EMAIL_HASH},
			getItemErrorToReturn:  nil,
			expectedErrorString:   "Can not create profile. Already exists",
			expectedResultString:  fmt.Sprintf("{\"ProfileID\":\"%s\"}\n", ""),
		},
		{
			description:           "Create failed - profile already exists",
			profile:               entity.ProfileData{Email: EMAIL},
			putItemOutputToReturn: nil,
			putItemErrorToReturn:  nil,
			getItemOutput:         nil,
			getItemReturnObject:   nil,
			getItemErrorToReturn:  errors.New("puke"),
			expectedErrorString:   "puke",
			expectedResultString:  fmt.Sprintf("{\"ProfileID\":\"%s\"}\n", ""),
		},
	}

	for _, tC := range testCases {
		logger := zerolog.New(os.Stdout)

		mockService := ServiceImpl{
			Client: dbclient.ClientImpl{
				Conn: mock.DB{
					PutItemOutputToReturn: tC.putItemOutputToReturn,
					PutItemErrorToReturn:  tC.putItemErrorToReturn,
					GetItemOutputToReturn: tC.getItemOutput,
					GetItemErrorToReturn:  tC.getItemErrorToReturn,
					GetItemReturnObject:   tC.getItemReturnObject,
				},
				Logger: &logger,
			},
			Logger: &logger,
		}

		result, err := mockService.CreateProfile(context.Background(), tC.profile)
		if tC.expectedErrorString != "" {
			assert.Equal(t, tC.expectedErrorString, err.Error())
		} else {
			assert.NoError(t, err)
		}

		buffer := new(bytes.Buffer)
		json.NewEncoder(buffer).Encode(result)
		assert.Equal(t, tC.expectedResultString, string(buffer.Bytes()))
	}
}
