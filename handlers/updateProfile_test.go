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

	"github.com/gorilla/mux"
	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/dbclient/mock"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

func TestUpdateProfile(t *testing.T) {
	testCases := []struct {
		description                string
		profileID                  string
		bodyToSend                 string
		expectedStatusCode         int
		expectedResponseBodyResult string
		PutItemErrorToReturn       error
		GetItemOutputToReturn      *dynamodb.GetItemOutput
		GetItemReturnObject        interface{}
		GetItemErrorToReturn       error
	}{
		{
			description:                "Happy path",
			profileID:                  EMAIL_HASH,
			bodyToSend:                 `{"email":"` + EMAIL + `"}`,
			expectedStatusCode:         200,
			expectedResponseBodyResult: "",
			PutItemErrorToReturn:       nil,
			GetItemOutputToReturn:      &dynamodb.GetItemOutput{},
			GetItemReturnObject:        entity.Profile{ID: EMAIL_HASH},
			GetItemErrorToReturn:       nil,
		},
		{
			description:                "Bad req body",
			profileID:                  EMAIL_HASH,
			bodyToSend:                 `rtr39gk402apg"}`,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"status\":\"Bad Request\",\"message\":\"Bad req body\",\"error\":\"invalid character 'r' looking for beginning of value\"}",
			PutItemErrorToReturn:       nil,
			GetItemOutputToReturn:      nil,
			GetItemReturnObject:        nil,
			GetItemErrorToReturn:       nil,
		},
		{
			description:                "Attempted to update email",
			profileID:                  EMAIL_HASH,
			bodyToSend:                 `{"email":"` + "appendAdditionalData" + EMAIL + `"}`,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"status\":\"Bad Request\",\"message\":\"UpdateProfile failed: attempted to change email\",\"error\":\"Email inconsistent with ProfileID\"}",
			PutItemErrorToReturn:       nil,
			GetItemOutputToReturn:      nil,
			GetItemReturnObject:        nil,
			GetItemErrorToReturn:       nil,
		},
		{
			description:                "Validation error - missing required attributes - Email",
			profileID:                  EMAIL_HASH,
			bodyToSend:                 `{"email":""}`,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"status\":\"Bad Request\",\"message\":\"Profile validation failed\",\"error\":\"Key: 'ProfileData.Email' Error:Field validation for 'Email' failed on the 'required' tag\"}",
			PutItemErrorToReturn:       nil,
			GetItemOutputToReturn:      nil,
			GetItemReturnObject:        nil,
			GetItemErrorToReturn:       nil,
		},
		{
			description:                "DB Error - PutItem failed",
			profileID:                  EMAIL_HASH,
			bodyToSend:                 `{"email":"` + EMAIL + `"}`,
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"status\":\"Internal Server Error\",\"message\":\"UpdateProfile failed\",\"error\":\"puke\"}",
			PutItemErrorToReturn:       errors.New("puke"),
			GetItemOutputToReturn:      &dynamodb.GetItemOutput{},
			GetItemReturnObject:        entity.Profile{ID: "123"},
			GetItemErrorToReturn:       nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			asserter := assert.New(t)
			logger := zerolog.New(os.Stdout)

			mockService := service.ServiceImpl{
				Client: dbclient.ClientImpl{
					DynamoDB: mock.DB{
						GetItemOutputToReturn: tC.GetItemOutputToReturn,
						GetItemReturnObject:   tC.GetItemReturnObject,
						GetItemErrorToReturn:  tC.GetItemErrorToReturn,
						PutItemErrorToReturn:  tC.PutItemErrorToReturn,
					},
					Logger: &logger,
				},
				Logger: &logger,
			}

			r := mux.NewRouter()
			s := r.PathPrefix("/test").Subrouter()
			s.HandleFunc("/{id}", UpdateProfile(mockService, validator.New())).Methods(http.MethodPut)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/test/%s", ts.URL, tC.profileID), bytes.NewReader([]byte(tC.bodyToSend)))

			res, err := ts.Client().Do(req)

			asserter.NoError(err)
			asserter.Equal(tC.expectedStatusCode, res.StatusCode)

			body, err := ioutil.ReadAll(res.Body)
			asserter.Equal(tC.expectedResponseBodyResult, string(body))
		})
	}
}
