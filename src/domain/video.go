package domain

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"github.com/justepl2/gopro_to_gpx_api/tools"
	"gorm.io/gorm"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusDone    Status = "done"
	StatusError   Status = "error"
)

type Video struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"column:name"`
	FilePath    string    `gorm:"column:file_path"`
	FileName    string    `gorm:"column:file_name"`
	FileType    string    `gorm:"column:file_type"`
	Duration    float64   `gorm:"column:duration"`
	CameraModel string    `gorm:"column:camera_model"`
	S3Location  string    `gorm:"column:s3_location"`
	Status      Status    `gorm:"column:status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (dv *Video) FromRequest(rv request.CreateVideo) {
	dv.Name = rv.Name
	dv.FilePath = rv.FilePath
}

// get Data from FilePath
func (dv *Video) FillVideoMetadata() error {
	// Use Exiftool CLI to get metadata, try to get Camera Model, Duration and FileName and FileType
	cmd := exec.Command("exiftool", "-j", "-Model", "-Duration", "-FileName", "-FileTypeExtension", dv.FilePath)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// parse output
	var data []map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return err
	}

	// Store information on Video struct
	if fileName, ok := data[0]["FileName"].(string); ok {
		dv.FileName = strings.Trim(fileName, ".MP4")
	}

	if fileType, ok := data[0]["FileTypeExtension"].(string); ok {
		dv.FileType = fileType
	}

	if model, ok := data[0]["Model"].(string); ok {
		dv.CameraModel = model
	}

	if duration, ok := data[0]["Duration"].(string); ok {
		duration = strings.Trim(duration, " s")
		duration, err := strconv.ParseFloat(duration, 32)
		if err != nil {
			return err
		}
		dv.Duration = tools.TruncateFloat(duration)
	}

	return nil
}

func (dv *Video) TransformMp4FileToBinFile() error {
	rawFileDestDir := os.Getenv("RAW_VIDEO_DEST_DIR")
	cmd := exec.Command("ffmpeg", "-y", "-i", dv.FilePath, "-codec", "copy", "-map", "0:3", "-f", "rawvideo", rawFileDestDir+dv.FileName+".bin")

	return cmd.Run()
}
