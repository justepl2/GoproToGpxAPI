package request

import "errors"

type RefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}

func (r *RefreshToken) Validate() error {
	if r.RefreshToken == "" {
		return errors.New("refresh token is required")
	}
	return nil
}
