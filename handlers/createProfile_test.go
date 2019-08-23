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
	"github.com/teohrt/cruddyAPI/dbclient/mock"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

const (
	EMAIL      = "foo@bar.com"
	EMAIL_HASH = "0c7e6a405862e402eb76a70f8a26fc732d07c32931e9fae9ab1582911d2e8a3b"
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
				"firstName": "Trace",
				"lastName": "Ohrt",
				"address": {
					"street": "175 Calvert Dr",
					"city": "Cupertino",
					"state": "California",
					"zipcode": "95014"
				},
				"email": "` + EMAIL + `"
				}`,
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{ID: ""},
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      &dynamodb.PutItemOutput{},
			putItemErrorToReturn:       nil,
			expectedStatusCode:         201,
			expectedResponseBodyResult: `{"ProfileID":"` + EMAIL_HASH + `"}`,
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
			expectedResponseBodyResult: "{\"status\":\"Bad Request\",\"message\":\"Bad req body\",\"error\":\"invalid character 'p' looking for beginning of value\"}",
		},
		{
			description:                "Req body validation failed - firstname is not alpha",
			bodyToSend:                 `{"firstName": "trace3", "email": "` + EMAIL + `"}`,
			getItemOutputToReturn:      nil,
			getItemReturnObject:        nil,
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      nil,
			putItemErrorToReturn:       nil,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"status\":\"Bad Request\",\"message\":\"Profile validation failed\",\"error\":\"Key: 'ProfileData.FirstName' Error:Field validation for 'FirstName' failed on the 'alpha' tag\"}",
		},
		{
			description:                "Req body validation failed - bad email",
			bodyToSend:                 `{"email": "foobar.com"}`,
			getItemOutputToReturn:      nil,
			getItemReturnObject:        nil,
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      nil,
			putItemErrorToReturn:       nil,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"status\":\"Bad Request\",\"message\":\"Profile validation failed\",\"error\":\"Key: 'ProfileData.Email' Error:Field validation for 'Email' failed on the 'email' tag\"}",
		},
		{
			description: "CreateProfile service fails - pukes on prexisting profile",
			bodyToSend: `{
				"firstName": "Trace",
				"lastName": "Ohrt",
				"address": {
					"street": "175 Calvert Dr",
					"city": "Cupertino",
					"state": "California",
					"zipcode": "95014"
				},
				"email": "` + EMAIL + `"
				}`,
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{ID: EMAIL_HASH},
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      &dynamodb.PutItemOutput{},
			putItemErrorToReturn:       nil,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"status\":\"Bad Request\",\"message\":\"Profile already exists\",\"error\":\"Can not create profile. Already exists\"}",
		},
		{
			description: "CreateProfile service fails - GetItem pukes",
			bodyToSend: `{
				"firstName": "Trace",
				"lastName": "Ohrt",
				"address": {
					"street": "175 Calvert Dr",
					"city": "Cupertino",
					"state": "California",
					"zipcode": "95014"
				},
				"email": "` + EMAIL + `"
				}`,
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        nil,
			getItemErrorToReturn:       errors.New("puke"),
			putItemOutputToReturn:      &dynamodb.PutItemOutput{},
			putItemErrorToReturn:       nil,
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"status\":\"Internal Server Error\",\"message\":\"Adding profile failed\",\"error\":\"puke\"}",
		},
		{
			description: "CreateProfile service fails - PutItem pukes",
			bodyToSend: `{
				"firstName": "Trace",
				"lastName": "Ohrt",
				"address": {
					"street": "175 Calvert Dr",
					"city": "Cupertino",
					"state": "California",
					"zipcode": "95014"
				},
				"email": "` + EMAIL + `"
				}`,
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{ID: ""},
			getItemErrorToReturn:       nil,
			putItemOutputToReturn:      nil,
			putItemErrorToReturn:       errors.New("puke"),
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"status\":\"Internal Server Error\",\"message\":\"Adding profile failed\",\"error\":\"puke\"}",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			asserter := assert.New(t)
			logger := zerolog.New(os.Stdout)

			mockService := service.ServiceImpl{
				Client: dbclient.ClientImpl{
					DynamoDB: mock.DB{
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
