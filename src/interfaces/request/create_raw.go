package request

import (
	"errors"

	"github.com/google/uuid"
)

type RawFile struct {
	Name   string    `json:"name" validate:"required" example:"video_1.bin"`
	File   []byte    `json:"file" validate:"required"`
	UserId uuid.UUID `json:"userId" example:"5f5e3e4e-3e4e-5f5e-3e4e-5f5e3e4e3e4e"`
}

func (rf *RawFile) Validate() error {
	if rf.Name == "" {
		return errors.New("name is required")
	}

	if len(rf.File) == 0 {
		return errors.New("file is required")
	}

	if rf.UserId == uuid.Nil {
		return errors.New("UserId isn't on jwt token, checkout your jwt")
	}

	return nil
}
