package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"
	"github.com/teohrt/cruddyAPI/utils"
	"gopkg.in/go-playground/validator.v9"
)

func UpdateProfile(svc service.Service, v *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		decoder := json.NewDecoder(r.Body)
		params := mux.Vars(r)
		profileID := params["id"]

		profileData := new(entity.ProfileData)
		if err := decoder.Decode(profileData); err != nil {
			logger.Debug().Err(err).Msg("Bad req body")
			utils.RespondWithError("Bad req body", err, http.StatusBadRequest, w)
			return
		}

		if err := validateProfile(profileData, v); err != nil {
			logger.Debug().Err(err).Msg("Profile validation failed")
			utils.RespondWithError("Profile validation failed", err, http.StatusBadRequest, w)
			return
		}

		if err := svc.UpdateProfile(r.Context(), *profileData, profileID); err != nil {
			logger.Error().Err(err).Msg("UpdateProfile failed")
			utils.RespondWithError("UpdateProfile failed", err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}
}
