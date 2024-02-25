package request

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email" example:"test@test.com"`
}

func (fp *ForgotPassword) Validate() error {
	validate := validator.New()
	err := validate.Struct(fp)
	if err != nil {
		return err
	}

	if fp.Email == "" {
		return errors.New("email is required")
	}
	return nil
}
