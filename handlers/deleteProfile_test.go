package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/dbclient/mock"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProfileHandler(t *testing.T) {
	testCases := []struct {
		description                string
		profileID                  string
		getItemOutputToReturn      *dynamodb.GetItemOutput
		getItemReturnObject        interface{}
		getItemErrorToReturn       error
		deleteItemErrorToReturn    error
		expectedStatusCode         int
		expectedResponseBodyResult string
	}{
		{
			description:                "Happy path",
			profileID:                  "123",
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{ID: "123"},
			getItemErrorToReturn:       nil,
			deleteItemErrorToReturn:    nil,
			expectedStatusCode:         200,
			expectedResponseBodyResult: "",
		},
		{
			description:                "Profile doesn't exist",
			profileID:                  "123",
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{},
			getItemErrorToReturn:       nil,
			deleteItemErrorToReturn:    nil,
			expectedStatusCode:         404,
			expectedResponseBodyResult: "{\"status\":\"Not Found\",\"message\":\"Profile not found\",\"error\":\"Could not find profile associated with: 123\"}",
		},
		{
			description:                "DB error - GetProfile puked",
			profileID:                  "123",
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        nil,
			getItemErrorToReturn:       errors.New("puke"),
			deleteItemErrorToReturn:    nil,
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"status\":\"Internal Server Error\",\"message\":\"DeleteProfile failed\",\"error\":\"puke\"}",
		},
		{
			description:                "DB error - DeleteProfile failed",
			profileID:                  "123",
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{ID: "123"},
			getItemErrorToReturn:       nil,
			deleteItemErrorToReturn:    errors.New("Puke"),
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"status\":\"Internal Server Error\",\"message\":\"DeleteProfile failed\",\"error\":\"Puke\"}",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			asserter := assert.New(t)
			logger := zerolog.New(os.Stdout)

			mockService := service.ServiceImpl{
				Client: dbclient.ClientImpl{
					DynamoDB: mock.DB{
						DeleteItemErrorToReturn: tC.deleteItemErrorToReturn,

						GetItemOutputToReturn: tC.getItemOutputToReturn,
						GetItemReturnObject:   tC.getItemReturnObject,
						GetItemErrorToReturn:  tC.getItemErrorToReturn,
					},
					Logger: &logger,
				},
				Logger: &logger,
			}

			r := mux.NewRouter()
			s := r.PathPrefix("/test").Subrouter()
			s.HandleFunc("/{id}", DeleteProfile(mockService)).Methods(http.MethodDelete)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/test/%s", ts.URL, tC.profileID), nil)
			res, err := ts.Client().Do(req)

			asserter.NoError(err)
			asserter.Equal(tC.expectedStatusCode, res.StatusCode)

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			asserter.Equal(tC.expectedResponseBodyResult, string(body))
		})
	}
}
