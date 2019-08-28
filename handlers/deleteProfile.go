package handlers

import (
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/service"
	"github.com/teohrt/cruddyAPI/utils"
)

func DeleteProfile(svc service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, seg := xray.BeginSegment(r.Context(), "DeleteProfile handler")
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		profileID := chi.URLParam(r, "id")

		if err := svc.DeleteProfile(ctx, profileID); err != nil {
			seg.Close(err)
			switch err.(type) {
			case service.ProfileNotFoundError:
				msg := "Profile not found"
				logger.Error().Err(err).Msg(msg)
				utils.RespondWithError(msg, err, http.StatusNotFound, w)
				return
			default:
				msg := "DeleteProfile failed"
				logger.Error().Err(err).Msg(msg)
				utils.RespondWithError(msg, err, http.StatusInternalServerError, w)
				return
			}
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		seg.Close(nil)
		return
	}
}
