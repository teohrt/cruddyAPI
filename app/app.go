package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"

	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/handlers"
	"github.com/teohrt/cruddyAPI/service"
)

func Start() {
	SERVER_PORT := os.Getenv("_LAMBDA_SERVER_PORT")

	config := dbclient.Config{}
	env.Parse(&config)
	svc := service.New(&config)
	v := validator.New()

	r := mux.NewRouter()
	s := r.PathPrefix("/cruddyAPI/v1").Subrouter()
	s.HandleFunc("/profiles", handlers.CreateProfile(svc, v)).Methods(http.MethodPost)
	s.HandleFunc("/profiles/{id}", handlers.GetProfile(svc)).Methods(http.MethodGet)
	s.HandleFunc("/profiles/{id}", handlers.UpdateProfile(svc, v)).Methods(http.MethodPut)
	s.HandleFunc("/profiles/{id}", handlers.DeleteProfile(svc)).Methods(http.MethodDelete)

	fmt.Println("Server listening on port :" + SERVER_PORT)
	log.Fatal(http.ListenAndServe(":"+SERVER_PORT, r))
}
