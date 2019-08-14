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

func TestUpdateProfile(t *testing.T) {
	testCases := []struct {
		description                string
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
			bodyToSend:                 `{"id": "123", "email":"foo@bar.com"}`,
			expectedStatusCode:         200,
			expectedResponseBodyResult: "",
			PutItemErrorToReturn:       nil,
			GetItemOutputToReturn:      &dynamodb.GetItemOutput{},
			GetItemReturnObject:        entity.Profile{ID: "123"},
			GetItemErrorToReturn:       nil,
		},
		{
			description:                "Bad req body",
			bodyToSend:                 `rtr39gk402apg"}`,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"Message\":\"Bad req body\",\"Error\":\"invalid character 'r' looking for beginning of value\"}",
			PutItemErrorToReturn:       nil,
			GetItemOutputToReturn:      nil,
			GetItemReturnObject:        nil,
			GetItemErrorToReturn:       nil,
		},
		{
			description:                "Validation error - missing required attributes - ID",
			bodyToSend:                 `{"email":"foo@bar.com"}`,
			expectedStatusCode:         400,
			expectedResponseBodyResult: "{\"Message\":\"Profile validation failed\",\"Error\":\"Key: 'Profile.ID' Error:Field validation for 'ID' failed on the 'required' tag\"}",
			PutItemErrorToReturn:       nil,
			GetItemOutputToReturn:      nil,
			GetItemReturnObject:        nil,
			GetItemErrorToReturn:       nil,
		},
		{
			description:                "DB Error - PutItem failed",
			bodyToSend:                 `{"id": "123", "email":"foo@bar.com"}`,
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"Message\":\"UpdateProfile failed\",\"Error\":\"puke\"}",
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
					DynamoDB: testutils.MockDB{
						GetItemOutputToReturn: tC.GetItemOutputToReturn,
						GetItemReturnObject:   tC.GetItemReturnObject,
						GetItemErrorToReturn:  tC.GetItemErrorToReturn,
						PutItemErrorToReturn:  tC.PutItemErrorToReturn,
					},
					Logger: &logger,
				},
				Logger: &logger,
			}

			r := chi.NewRouter()
			r.Put("/test", UpdateProfile(mockService, validator.New()))
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/test", ts.URL), bytes.NewReader([]byte(tC.bodyToSend)))
			res, err := ts.Client().Do(req)

			asserter.NoError(err)
			asserter.Equal(tC.expectedStatusCode, res.StatusCode)

			body, err := ioutil.ReadAll(res.Body)
			asserter.Equal(tC.expectedResponseBodyResult, string(body))
		})
	}
}
