package request

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type MailValidation struct {
	Username         string `json:"username" validate:"required" example:"john_doe"`
	ConfirmationCode string `json:"confirmationCode" validate:"required" example:"123456"`
}

func (mv *MailValidation) Validate() error {
	validate := validator.New()
	err := validate.Struct(mv)
	if err != nil {
		return err
	}

	if mv.Username == "" {
		return errors.New("username is required")
	}

	if mv.ConfirmationCode == "" {
		return errors.New("confirmationCode is required")
	}

	if len(mv.ConfirmationCode) != 6 {
		return errors.New("confirmationCode must be 6 characters long")
	}

	return nil
}
