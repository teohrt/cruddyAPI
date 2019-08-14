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
	"github.com/teohrt/cruddyAPI/service"
	"github.com/teohrt/cruddyAPI/testutils"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProfileHandler(t *testing.T) {
	testCases := []struct {
		description                string
		profileID                  string
		deleteItemErrorToReturn    error
		expectedStatusCode         int
		expectedResponseBodyResult string
	}{
		{
			description:                "Happy path",
			profileID:                  "123",
			deleteItemErrorToReturn:    nil,
			expectedStatusCode:         200,
			expectedResponseBodyResult: "",
		},
		{
			description:                "DB error - DeleteProfile failed",
			profileID:                  "123",
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
					DynamoDB: testutils.MockDB{
						DeleteItemErrorToReturn: tC.deleteItemErrorToReturn,
					},
					Logger: &logger,
				},
				Logger: &logger,
			}

			r := chi.NewRouter()
			r.Delete("/test/{id}", DeleteProfile(mockService))
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
