package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"
	"github.com/teohrt/cruddyAPI/testutils"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

func TestCreateProfileHandler(t *testing.T) {
	testCases := []struct {
		description string
		bodyToSend  string

		getItemOutputToReturn *dynamodb.GetItemOutput
		getItemReturnObject   interface{}
		getItemErrorToReturn  error

		putItemOutputToReturn *dynamodb.PutItemOutput
		putItemErrorToReturn  error

		expectedStatusCode         int
		expectedResponseBodyResult string
	}{
		{
			description: "Happy path",
			bodyToSend: `{
				"id": "123",
				"firstName": "Trace",
				"lastName": "Ohrt",
				"address": {
					"street": "175 Calvert Dr",
					"city": "Cupertino",
					"state": "California",
					"zipcode": "95014"
				},
				"email": "teohrt18@gmail.com"
			}`,
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{ID: ""},
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      &dynamodb.PutItemOutput{},
			putItemErrorToReturn:       nil,
			expectedStatusCode:         201,
			expectedResponseBodyResult: `{"ProfileID":"123"}`,
		},
		{
			description:                "Bady req body",
			bodyToSend:                 `puke`,
			getItemOutputToReturn:      nil,
			getItemReturnObject:        nil,
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      nil,
			putItemErrorToReturn:       nil,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"Message\":\"Bad req body\",\"Error\":\"invalid character 'p' looking for beginning of value\"}",
		},
		{
			description:                "Req body validation failed - firstname is not alpha",
			bodyToSend:                 `{"ID": "123", "firstName": "trace3"}`,
			getItemOutputToReturn:      nil,
			getItemReturnObject:        nil,
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      nil,
			putItemErrorToReturn:       nil,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"Message\":\"Profile validation failed\",\"Error\":\"Key: 'Profile.FirstName' Error:Field validation for 'FirstName' failed on the 'alpha' tag\"}",
		},
		{
			description:                "Req body validation failed - bad email",
			bodyToSend:                 `{"ID": "123", "email": "foobar.com"}`,
			getItemOutputToReturn:      nil,
			getItemReturnObject:        nil,
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      nil,
			putItemErrorToReturn:       nil,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"Message\":\"Profile validation failed\",\"Error\":\"Key: 'Profile.Email' Error:Field validation for 'Email' failed on the 'email' tag\"}",
		},
		{
			description: "CreateProfile service fails - pukes on prexisting profile",
			bodyToSend: `{
				"id": "123",
				"firstName": "Trace",
				"lastName": "Ohrt",
				"address": {
					"street": "175 Calvert Dr",
					"city": "Cupertino",
					"state": "California",
					"zipcode": "95014"
				},
				"email": "teohrt18@gmail.com"
			}`,
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{ID: "123"},
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      &dynamodb.PutItemOutput{},
			putItemErrorToReturn:       nil,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"Message\":\"Profile already exists\",\"Error\":\"Can not create profile. Already exists\"}",
		},
		{
			description: "CreateProfile service fails - GetItem pukes",
			bodyToSend: `{
				"id": "123",
				"firstName": "Trace",
				"lastName": "Ohrt",
				"address": {
					"street": "175 Calvert Dr",
					"city": "Cupertino",
					"state": "California",
					"zipcode": "95014"
				},
				"email": "teohrt18@gmail.com"
			}`,
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        nil,
			getItemErrorToReturn:       errors.New("puke"),
			putItemOutputToReturn:      &dynamodb.PutItemOutput{},
			putItemErrorToReturn:       nil,
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"Message\":\"Adding profile failed\",\"Error\":\"puke\"}",
		},
		{
			description: "CreateProfile service fails - PutItem pukes",
			bodyToSend: `{
				"id": "123",
				"firstName": "Trace",
				"lastName": "Ohrt",
				"address": {
					"street": "175 Calvert Dr",
					"city": "Cupertino",
					"state": "California",
					"zipcode": "95014"
				},
				"email": "teohrt18@gmail.com"
			}`,
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{ID: ""},
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      nil,
			putItemErrorToReturn:       errors.New("puke"),
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"Message\":\"Adding profile failed\",\"Error\":\"puke\"}",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			asserter := assert.New(t)
			logger := zerolog.New(os.Stdout)

			mockService := service.ServiceImpl{
				Client: dbclient.ClientImpl{
					DynamoDB: testutils.MockDB{
						GetItemOutputToReturn: tC.getItemOutputToReturn,
						GetItemReturnObject:   tC.getItemReturnObject,
						GetItemErrorToReturn:  tC.getItemErrorToReturn,

						PutItemOutputToReturn: tC.putItemOutputToReturn,
						PutItemErrorToReturn:  tC.putItemErrorToReturn,
					},
					Logger: &logger,
				},
				Logger: &logger,
			}

			r := chi.NewRouter()
			r.Post("/test", CreateProfile(mockService, validator.New()))
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest("POST", fmt.Sprintf("%s/test", ts.URL), bytes.NewReader([]byte(tC.bodyToSend)))
			res, err := ts.Client().Do(req)

			asserter.NoError(err)
			asserter.Equal(tC.expectedStatusCode, res.StatusCode)

			body, err := ioutil.ReadAll(res.Body)
			asserter.Equal(tC.expectedResponseBodyResult, string(body))
		})
	}
}
