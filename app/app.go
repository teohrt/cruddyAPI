package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/handlers"
	"github.com/teohrt/cruddyAPI/service"

	"github.com/caarlos0/env"
)

type Config struct {
	Port string
}

func Start(c Config) {
	dbconfig := dbclient.DBConfig{}
	env.Parse(&dbconfig)
	svc := service.New(&dbconfig)

	r := mux.NewRouter()
	s := r.PathPrefix("/cruddyAPI").Subrouter()
	s.HandleFunc("/profile/{id}", handlers.GetProfileHandler(svc)).Methods("GET")
	s.HandleFunc("/profile", handlers.CreateProfileHandler(svc)).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+c.Port, r))
}
