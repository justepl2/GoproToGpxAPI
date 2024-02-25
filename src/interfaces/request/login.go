package request

import "errors"

type Login struct {
	Email    string `json:"email" validate:"required,email" example:"test@test.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

func (cu *Login) Validate() error {
	if cu.Email == "" {
		return errors.New("email is required")
	}

	if cu.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
