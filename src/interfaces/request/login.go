package request

import "errors"

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
