package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"gopkg.in/go-playground/validator.v9"

	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"
)

// TODO
func CreateProfile(svc service.Service, v *validator.Validate) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger := zerolog.Ctx(r.Context())
		decoder := json.NewDecoder(r.Body)

		profile := new(entity.Profile)
		if err := decoder.Decode(profile); err != nil {
			logger.Debug().Err(err).Msg("Bad req body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := validateProfile(profile, v); err != nil {
			logger.Debug().Err(err).Msg("Profile validation failed")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := svc.CreateProfile(r.Context(), *profile)
		if err != nil {
			logger.Error().Err(err).Msg("Adding profile failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonObj, err := json.Marshal(result)
		if err != nil {
			logger.Warn().Err(err).Msg("Failed marshalling json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", fmt.Sprintf("/profile/id/%v", result.ProfileID))
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonObj)
		return
	}
}
