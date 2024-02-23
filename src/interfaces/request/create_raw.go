package request

import (
	"errors"

	"github.com/google/uuid"
)

type RawFile struct {
	Name   string
	File   []byte
	UserId uuid.UUID
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
