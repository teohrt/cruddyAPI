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

		profile := new(entity.ProfileData)
		if err := decoder.Decode(profile); err != nil {
			msg := "Bad req body"
			logger.Debug().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		if err := validateProfile(profile, v); err != nil {
			msg := "Profile validation failed"
			logger.Debug().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		result, err := svc.CreateProfile(r.Context(), *profile)
		if err != nil {
			switch err.(type) {
			case service.ProfileAlreadyExistsError:
				w.Header().Add("Location", fmt.Sprintf("api/v1/profiles/%v", utils.Hash(profile.Email)))
				msg := "Profile already exists"
				logger.Debug().Err(err).Msg(msg)
				utils.RespondWithError(msg, err, http.StatusBadRequest, w)
				return
			default:
				msg := "Adding profile failed"
				logger.Error().Err(err).Msg(msg)
				utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
				return
			}
		}

		jsonObj, err := json.Marshal(result)
		if err != nil {
			msg := "Failed marshalling json"
			logger.Warn().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Location", fmt.Sprintf("api/v1/profiles/%v", result.ProfileID))
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonObj)
		return
	}
}
