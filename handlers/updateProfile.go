package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"
	"gopkg.in/go-playground/validator.v9"
)

// TODO
func UpdateProfile(svc service.Service, v *validator.Validate) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
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

		fmt.Fprintf(w, "TODO")
	}
}
