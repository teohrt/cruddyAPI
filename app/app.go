package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teohrt/cruddyAPI/handlers"
)

type Config struct {
	Port string
}

func Start(c Config) {
	r := mux.NewRouter()

	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/echo/{word}", handlers.EchoHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+c.Port, r))
}
