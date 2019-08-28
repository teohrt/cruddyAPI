package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/entity"
	"github.com/teohrt/cruddyAPI/service"
	"github.com/teohrt/cruddyAPI/utils"
	"gopkg.in/go-playground/validator.v9"
)

func UpdateProfile(svc service.Service, v *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		ctx, seg := utils.StartXraySegment(r.Context(), "UpdateProfile handler")
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		decoder := json.NewDecoder(r.Body)
		profileID := chi.URLParam(r, "id")

		profileData := new(entity.ProfileData)
		if err := decoder.Decode(profileData); err != nil {
			seg.Close(err)
			msg := "Bad req body"
			logger.Debug().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		if err := validateProfile(profileData, v); err != nil {
			seg.Close(err)
			msg := "Profile validation failed"
			logger.Debug().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusBadRequest, w)
			return
		}

		if err := svc.UpdateProfile(ctx, *profileData, profileID); err != nil {
			if err != nil {
				seg.Close(err)
				switch err.(type) {
				case service.EmailIncsonsistentWithProfileIDError:
					msg := "UpdateProfile failed: attempted to change email"
					logger.Error().Err(err).Msg(msg)
					utils.RespondWithError(msg, err, http.StatusBadRequest, w)
					return
				default:
					msg := "UpdateProfile failed"
					logger.Error().Err(err).Msg(msg)
					utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
					return
				}
			}
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		seg.Close(nil)
		return
	}
}
