package request

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CreateVideo struct {
	Name     string `json:"name" validate:"required" example:"video_1.mp4"`
	FilePath string `json:"filePath" validate:"required" example:"/path/to/video_1.mp4"`
}

func (cv *CreateVideo) Validate() error {
	validate := validator.New()
	err := validate.Struct(cv)
	if err != nil {
		return err
	}

	if cv.Name == "" {
		return errors.New("name is required")
	}

	if cv.FilePath == "" {
		return errors.New("filePath is required")
	}

	file, err := os.Open(cv.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use MIME type to check if the file is a video
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	mimeType := http.DetectContentType(buffer)

	// Check if the file is a video
	if !strings.HasPrefix(mimeType, "video") {
		return errors.New("filePath must be a video file")
	}

	return nil
}
