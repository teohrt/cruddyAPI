package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teohrt/cruddyAPI/service"
)

func GetProfileHandler(svc service.Service) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]
		fmt.Fprintf(w, word)
	}
}
