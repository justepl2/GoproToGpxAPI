package request

import "errors"

type Signup struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func (cu *Signup) Validate() error {
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
