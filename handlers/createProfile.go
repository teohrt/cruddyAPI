package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"

	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"
	"github.com/teohrt/cruddyAPI/utils"
)

func CreateProfile(svc service.Service, v *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		decoder := json.NewDecoder(r.Body)

		profile := new(entity.Profile)
		if err := decoder.Decode(profile); err != nil {
			logger.Debug().Err(err).Msg("Bad req body")
			utils.RespondWithError("Bad req body", err, http.StatusBadRequest, w)
			return
		}

		if err := validateProfile(profile, v); err != nil {
			logger.Debug().Err(err).Msg("Profile validation failed")
			utils.RespondWithError("Profile validation failed", err, http.StatusBadRequest, w)
			return
		}

		result, err := svc.CreateProfile(r.Context(), *profile)
		if err != nil {
			switch err.(type) {
			case service.ProfileAlreadyExistsError:
				logger.Debug().Err(err).Msg("Profile already exists")
				utils.RespondWithError("Profile already exists", err, http.StatusBadRequest, w)
				return
			default:
				logger.Error().Err(err).Msg("Adding profile failed")
				utils.RespondWithError("Adding profile failed", err, http.StatusInternalServerError, w)
				return
			}
		}

		jsonObj, err := json.Marshal(result)
		if err != nil {
			logger.Warn().Err(err).Msg("Failed marshalling json")
			utils.RespondWithError("Failed marshalling json", err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Location", fmt.Sprintf("/profiles/%v", result.ProfileID))
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonObj)
		return
	}
}
