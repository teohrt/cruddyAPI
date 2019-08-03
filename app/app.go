package app

import (
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"

	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/handlers"
	"github.com/teohrt/cruddyAPI/service"
)

type Config struct {
	Port string
}

func Start(c Config) {
	dbconfig := dbclient.DBConfig{}
	env.Parse(&dbconfig)
	svc := service.New(&dbconfig)
	v := validator.New()

	r := mux.NewRouter()
	s := r.PathPrefix("/cruddyAPI").Subrouter()
	s.HandleFunc("/profile", handlers.CreateProfileHandler(svc, v)).Methods("POST")
	s.HandleFunc("/profile/{id}", handlers.GetProfileHandler(svc)).Methods("GET")
	s.HandleFunc("/profile/{id}", handlers.UpdateProfileHandler(svc, v)).Methods("PUT")
	s.HandleFunc("/profile/{id}", handlers.DeleteProfileHandler(svc)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+c.Port, r))
}
