package handlers

import (
	"net/http"

	"github.com/teohrt/cruddyAPI/entity"
	"gopkg.in/go-playground/validator.v9"
)

type Handler func(w http.ResponseWriter, r *http.Request)

func validateProfile(p *entity.Profile, v *validator.Validate) error {
	return v.Struct(p)
}
