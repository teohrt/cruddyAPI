package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/service"
	"github.com/teohrt/cruddyAPI/utils"
)

func GetProfile(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		params := mux.Vars(r)
		profileID := params["id"]

		result, err := svc.GetProfile(r.Context(), profileID)
		if err != nil {
			logger.Error().Err(err).Msg("Get profile failed")
			utils.RespondWithError("Get profile failed", err, http.StatusInternalServerError, w)
			return
		}

		jsonObj, err := json.Marshal(result)
		if err != nil {
			logger.Warn().Err(err).Msg("Failed marshalling json")
			utils.RespondWithError("Failed marshalling json", err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonObj)
		return
	}
}
