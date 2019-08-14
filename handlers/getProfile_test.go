package handlers

import (
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
)

func TestGetProfileHandler(t *testing.T) {
	testCases := []struct {
		description string
		profileID   string

		getItemOutputToReturn *dynamodb.GetItemOutput
		getItemReturnObject   interface{}
		getItemErrorToReturn  error

		expectedStatusCode         int
		expectedResponseBodyResult string
	}{
		{
			description:           "Happy path",
			profileID:             "123",
			getItemOutputToReturn: &dynamodb.GetItemOutput{},
			getItemReturnObject: entity.Profile{
				ID:        "123",
				FirstName: "Trace",
				LastName:  "Ohrt",
				Address: entity.Address{
					Street:  "175 Calvert Dr",
					City:    "Cupertino",
					State:   "CA",
					ZipCode: "95014",
				},
				Email: "teohrt18@gmail.com",
			},
			getItemErrorToReturn:       nil,
			expectedStatusCode:         200,
			expectedResponseBodyResult: "{\"id\":\"123\",\"firstName\":\"Trace\",\"lastName\":\"Ohrt\",\"address\":{\"street\":\"175 Calvert Dr\",\"city\":\"Cupertino\",\"state\":\"CA\",\"zipCode\":\"95014\"},\"email\":\"teohrt18@gmail.com\"}",
		},
		{
			description:                "DB Error - Profile doesn't exist",
			profileID:                  "123",
			getItemOutputToReturn:      nil,
			getItemReturnObject:        nil,
			getItemErrorToReturn:       errors.New("puke"),
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"status\":\"Internal Server Error\",\"message\":\"Get profile failed\",\"error\":\"puke\"}",
		},
		{
			description:                "DB Error - GetProfile failed",
			profileID:                  "123",
			getItemOutputToReturn:      &dynamodb.GetItemOutput{},
			getItemReturnObject:        entity.Profile{},
			getItemErrorToReturn:       nil,
			expectedStatusCode:         500,
			expectedResponseBodyResult: "{\"status\":\"Internal Server Error\",\"message\":\"Get profile failed\",\"error\":\"Could not find profile associated with: \"}",
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
					},
					Logger: &logger,
				},
				Logger: &logger,
			}

			r := chi.NewRouter()
			r.Get("/test/{id}", GetProfile(mockService))
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/test/%s", ts.URL, tC.profileID), nil)
			res, err := ts.Client().Do(req)

			asserter.NoError(err)
			asserter.Equal(tC.expectedStatusCode, res.StatusCode)

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			asserter.Equal(tC.expectedResponseBodyResult, string(body))
		})
	}
}
