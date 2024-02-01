package request

import "errors"

type CreateVideo struct {
	Name     string  `json:"name"`
	Duration float64 `json:"duration"`
}

func (cv *CreateVideo) Validate() error {
	if cv.Name == "" {
		return errors.New("name is required")
	}

	if cv.Duration <= 0 {
		return errors.New("duration must be greater than 0")
	}

	return nil
}
