package handlers

import (
	"github.com/teohrt/cruddyAPI/entity"
	"gopkg.in/go-playground/validator.v9"
)

func validateProfile(p *entity.ProfileData, v *validator.Validate) error {
	return v.Struct(p)
}
