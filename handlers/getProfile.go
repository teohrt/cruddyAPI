package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/service"
	"github.com/teohrt/cruddyAPI/utils"
)

func GetProfile(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, seg := xray.BeginSegment(r.Context(), "GetProfile handler")
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		profileID := chi.URLParam(r, "id")

		result, err := svc.GetProfile(ctx, profileID)
		if err != nil {
			seg.Close(err)
			switch err.(type) {
			case service.ProfileNotFoundError:
				msg := "Profile not found"
				logger.Error().Err(err).Msg(msg)
				utils.RespondWithError(msg, err, http.StatusNotFound, w)
				return
			default:
				msg := "Get profile failed"
				logger.Error().Err(err).Msg(msg)
				utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
				return
			}
		}

		jsonObj, err := json.Marshal(result)
		if err != nil {
			seg.Close(err)
			msg := "Failed marshalling json"
			logger.Warn().Err(err).Msg(msg)
			utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonObj)
		seg.Close(nil)
		return
	}
}
