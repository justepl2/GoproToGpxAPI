package request

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Login struct {
	Email    string `json:"email" validate:"required,email" example:"test@test.com"`
	Password string `json:"password" validate:"required" example:"Password123."`
}

func (cu *Login) Validate() error {
	validate := validator.New()
	err := validate.Struct(cu)
	if err != nil {
		return err
	}

	if cu.Email == "" {
		return errors.New("email is required")
	}

	if cu.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
