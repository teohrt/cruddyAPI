package dbclient

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/teohrt/cruddyAPI/testutils"
)

func TestNewDBClient(t *testing.T) {
	logger := zerolog.New(os.Stdout)
	client := New(&Config{
		DynamoDBTableName: "profiles",
		AWSRegion:         "us-east-2",
	}, &logger)

	assert.IsType(t, ClientImpl{}, client)
}
func TestGetItem(t *testing.T) {
	testCases := []struct {
		description           string
		valueName             string
		value                 string
		expectedAv            *map[string]*dynamodb.AttributeValue
		expectedErrorString   string
		getItemOutputToReturn *dynamodb.GetItemOutput
		getItemErrorToReturn  error
	}{
		{
			description:           "Happy Path",
			valueName:             "ProfileID",
			value:                 "1",
			expectedAv:            &map[string]*dynamodb.AttributeValue{},
			expectedErrorString:   "",
			getItemOutputToReturn: new(dynamodb.GetItemOutput),
			getItemErrorToReturn:  nil,
		},
		{
			description:           "DB connection/Query failed.",
			valueName:             "ProfileID",
			value:                 "1",
			expectedAv:            nil,
			expectedErrorString:   "puke",
			getItemOutputToReturn: nil,
			getItemErrorToReturn:  errors.New("puke"),
		},
	}

	for _, tC := range testCases {
		logger := zerolog.New(os.Stdout)

		mockClient := ClientImpl{
			DynamoDB: testutils.MockDB{
				GetItemOutputToReturn: tC.getItemOutputToReturn,
				GetItemErrorToReturn:  tC.getItemErrorToReturn,
			},
			Logger: &logger,
		}

		av, err := mockClient.GetItem(context.Background(), tC.valueName, tC.value)
		if tC.expectedErrorString != "" {
			assert.Equal(t, tC.expectedErrorString, err.Error())
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tC.expectedAv, av)
	}
}

func TestUpsertItem(t *testing.T) {
	testCases := []struct {
		description         string
		in                  interface{}
		expectedItem        *dynamodb.PutItemOutput
		expectedErrorString string
		queryOutputToReturn *dynamodb.PutItemOutput
		queryErrorToReturn  error
	}{
		{
			description:         "Happy path",
			expectedItem:        &dynamodb.PutItemOutput{},
			expectedErrorString: "",
			queryOutputToReturn: &dynamodb.PutItemOutput{},
			queryErrorToReturn:  nil,
		},
		{
			description:         "DB connection/Query failed.",
			expectedItem:        nil,
			expectedErrorString: "puke",
			queryOutputToReturn: nil,
			queryErrorToReturn:  errors.New("puke"),
		},
	}

	for _, tC := range testCases {
		logger := zerolog.New(os.Stdout)

		clientImpl := ClientImpl{
			DynamoDB: testutils.MockDB{
				PutItemOutputToReturn: tC.queryOutputToReturn,
				PutItemErrorToReturn:  tC.queryErrorToReturn,
			},
			Logger: &logger,
		}

		item, err := clientImpl.UpsertItem(context.Background(), tC.in)
		if tC.expectedErrorString != "" {
			assert.Equal(t, tC.expectedErrorString, err.Error())
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tC.expectedItem, item)
	}
}

func TestDeleteItem(t *testing.T) {
	testCases := []struct {
		description              string
		valueName                string
		value                    string
		expectedOutput           *dynamodb.DeleteItemOutput
		expectedErrorString      string
		deleteItemOutputToReturn *dynamodb.DeleteItemOutput
		deleteItemErrorToReturn  error
	}{
		{
			description:              "Happy Path",
			valueName:                "id",
			value:                    "1",
			expectedOutput:           new(dynamodb.DeleteItemOutput),
			expectedErrorString:      "",
			deleteItemOutputToReturn: new(dynamodb.DeleteItemOutput),
			deleteItemErrorToReturn:  nil,
		},
		{
			description:              "DB connection/Query failed.",
			valueName:                "id",
			value:                    "1",
			expectedOutput:           nil,
			expectedErrorString:      "puke",
			deleteItemOutputToReturn: nil,
			deleteItemErrorToReturn:  errors.New("puke"),
		},
	}

	for _, tC := range testCases {
		logger := zerolog.New(os.Stdout)

		mockClient := ClientImpl{
			DynamoDB: testutils.MockDB{
				DeleteItemOutputToReturn: tC.deleteItemOutputToReturn,
				DeleteItemErrorToReturn:  tC.deleteItemErrorToReturn,
			},
			Logger: &logger,
		}

		out, err := mockClient.DeleteItem(context.Background(), tC.valueName, tC.value)
		if tC.expectedErrorString != "" {
			assert.Equal(t, tC.expectedErrorString, err.Error())
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tC.expectedOutput, out)
	}
}
