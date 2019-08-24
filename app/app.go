package app

import (
	"log"

	"github.com/apex/gateway"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	validatorV9 "gopkg.in/go-playground/validator.v9"

	"github.com/teohrt/cruddyAPI/dbclient"
	"github.com/teohrt/cruddyAPI/handlers"
	"github.com/teohrt/cruddyAPI/service"
)

func Start() {
	// PORT := os.Getenv("SERVER_PORT")
	config := dbclient.Config{}
	env.Parse(&config)
	svc := service.New(&config)
	validator := validatorV9.New()

	router := chi.NewRouter()
	router.Route("/v1", func(subRouter chi.Router) {
		subRouter.Get("/health", handlers.Health())
		subRouter.Route("/profiles", func(r chi.Router) {
			r.Post("/", handlers.CreateProfile(svc, validator))
			r.Get("/{id}", handlers.GetProfile(svc))
			r.Put("/{id}", handlers.UpdateProfile(svc, validator))
			r.Delete("/{id}", handlers.DeleteProfile(svc))
		})
	})

	// fmt.Println("Server running locally and listening on port :" + PORT)
	// log.Fatal(http.ListenAndServe(":"+PORT, router))

	log.Fatal(gateway.ListenAndServe("", router))
}
