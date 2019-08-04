package handlers

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/teohrt/cruddyAPI/service"
)

func DeleteProfile(svc service.Service) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		params := mux.Vars(r)
		profileID := params["id"]

		if err := svc.DeleteProfile(r.Context(), profileID); err != nil {
			logger.Error().Err(err).Msg("DeleteProfile failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}
