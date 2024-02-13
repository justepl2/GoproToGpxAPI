package request

import (
	"errors"
)

type RawFile struct {
	Name string
	File []byte
}

func (rf *RawFile) Validate() error {
	if rf.Name == "" {
		return errors.New("name is required")
	}

	if len(rf.File) == 0 {
		return errors.New("file is required")
	}

	return nil
}
