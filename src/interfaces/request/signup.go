package request

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Signup struct {
	Email       string `json:"email" validate:"required,email" example:"john@doe.com"`
	Username    string `json:"username" validate:"required" example:"john_doe"`
	Password    string `json:"password" validate:"required" example:"password123"`
	FirstName   string `json:"firstname" example:"John"`
	LastName    string `json:"lastname" example:"Doe"`
	PhoneNumber string `json:"phonenumber" validate:"e164" example:"+1234567890"`
}

func (cu *Signup) Validate() error {
	validate := validator.New()
	err := validate.Struct(cu)
	if err != nil {
		return err
	}

	if cu.Email == "" {
		return errors.New("email is required")
	}

	if cu.Username == "" {
		return errors.New("username is required")
	}

	if cu.Password == "" {
		return errors.New("password is required")
	}

	if cu.FirstName == "" {
		return errors.New("firstname is required")
	}

	if cu.LastName == "" {
		return errors.New("lastname is required")
	}

	return nil
}
