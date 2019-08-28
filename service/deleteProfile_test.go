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

func TestDeleteProfile(t *testing.T) {
	testCases := []struct {
		description              string
		id                       string
		getItemOutputToReturn    *dynamodb.GetItemOutput
		getItemReturnObject      interface{}
		getItemErrorToReturn     error
		deleteItemOutputToReturn *dynamodb.DeleteItemOutput
		deleteItemErrorToReturn  error
		expectedErrorString      string
	}{
		{
			description:              "Happy path",
			id:                       "123",
			getItemOutputToReturn:    &dynamodb.GetItemOutput{},
			getItemReturnObject:      entity.Profile{ID: "123"},
			getItemErrorToReturn:     nil,
			deleteItemOutputToReturn: &dynamodb.DeleteItemOutput{},
			deleteItemErrorToReturn:  nil,
			expectedErrorString:      "",
		},
		{
			description:              "DB Error- GetProfile found no profile",
			id:                       "123",
			getItemOutputToReturn:    &dynamodb.GetItemOutput{},
			getItemReturnObject:      entity.Profile{ID: ""},
			getItemErrorToReturn:     nil,
			deleteItemOutputToReturn: &dynamodb.DeleteItemOutput{},
			deleteItemErrorToReturn:  nil,
			expectedErrorString:      "Could not find profile associated with: 123",
		},
		{
			description:              "DB Error- GetProfile puked",
			id:                       "123",
			getItemOutputToReturn:    &dynamodb.GetItemOutput{},
			getItemReturnObject:      entity.Profile{},
			getItemErrorToReturn:     errors.New("puke"),
			deleteItemOutputToReturn: &dynamodb.DeleteItemOutput{},
			deleteItemErrorToReturn:  nil,
			expectedErrorString:      "puke",
		},
		{
			description:              "DB Error- DeleteItem puked",
			id:                       "123",
			getItemOutputToReturn:    &dynamodb.GetItemOutput{},
			getItemReturnObject:      entity.Profile{ID: "123"},
			getItemErrorToReturn:     nil,
			deleteItemOutputToReturn: &dynamodb.DeleteItemOutput{},
			deleteItemErrorToReturn:  errors.New("puke"),
			expectedErrorString:      "puke",
		},
	}

	for _, tC := range testCases {
		logger := zerolog.New(os.Stdout)

		mockService := ServiceImpl{
			Client: dbclient.ClientImpl{
				Conn: mock.DB{
					GetItemOutputToReturn: tC.getItemOutputToReturn,
					GetItemReturnObject:   tC.getItemReturnObject,
					GetItemErrorToReturn:  tC.getItemErrorToReturn,

					DeleteItemOutputToReturn: tC.deleteItemOutputToReturn,
					DeleteItemErrorToReturn:  tC.deleteItemErrorToReturn,
				},
				Logger: &logger,
			},
			Logger: &logger,
		}

		err := mockService.DeleteProfile(context.Background(), tC.id)

		if tC.expectedErrorString != "" {
			assert.Error(t, err)
			assert.Equal(t, tC.expectedErrorString, err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
