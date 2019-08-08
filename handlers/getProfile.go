package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/service"
)

func GetProfile(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		params := mux.Vars(r)
		profileID := params["id"]

		result, err := svc.GetProfile(r.Context(), profileID)
		if err != nil {
			logger.Error().Err(err).Msg("Get profile failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonObj, err := json.Marshal(result)
		if err != nil {
			logger.Warn().Err(err).Msg("Failed marshalling json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonObj)
		return
	}
}
