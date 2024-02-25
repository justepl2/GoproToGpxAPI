package request

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ResetPassword struct {
	Email           string `json:"email" validate:"required,email" example:"`
	NewPassword     string `json:"newPassword" validate:"required" example:"Password123."`
	ConfirmationKey string `json:"confirmationKey" validate:"required" example:"123456"`
}

func (rp *ResetPassword) Validate() error {
	validate := validator.New()
	err := validate.Struct(rp)
	if err != nil {
		return err
	}

	if rp.Email == "" {
		return errors.New("email is required")
	}

	if rp.NewPassword == "" {
		return errors.New("newPassword is required")
	}

	if rp.ConfirmationKey == "" {
		return errors.New("confirmationKey is required")
	}

	if len(rp.ConfirmationKey) != 6 {
		return errors.New("confirmationKey must be 6 characters long")
	}

	return nil
}
