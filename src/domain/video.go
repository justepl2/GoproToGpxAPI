package domain

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name               string    `gorm:"column:name"`
	FilePath           string    `gorm:"column:file_path"`
	FileName           string    `gorm:"column:file_name"`
	FileType           string    `gorm:"column:file_type"`
	Duration           float64   `gorm:"column:duration"`
	CameraModel        string    `gorm:"column:camera_model"`
	MediaUniqueID      string    `gorm:"column:media_unique_id"`
	CameraSerialNumber string    `gorm:"column:camera_serial_number"`
	S3Location         string    `gorm:"column:s3_location"`
	Status             Status    `gorm:"column:status"`
	GpxID              uuid.UUID
	Gpx                Gpx            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (dv *Video) FromRequest(rv request.CreateVideo) {
	dv.Name = rv.Name
	dv.FilePath = rv.FilePath
}

func (dv *Video) FromRawRequest(rf request.RawFile) {
	dv.Name = strings.TrimSuffix(rf.Name, filepath.Ext(rf.Name))
	dv.FileName = rf.Name
	dv.FilePath = os.Getenv("RAW_VIDEO_DEST_DIR") + rf.Name
	dv.FileType = "bin"
}

// get Data from FilePath
func (dv *Video) FillVideoMetadata() error {
	// Use Exiftool CLI to get metadata, try to get Camera Model, Duration and FileName and FileType + MediaUniqueID + CameraSerialNumber
	cmd := exec.Command("exiftool", "-j", "-Model", "-Duration", "-FileName", "-FileTypeExtension", "-MediaUniqueID", "-CameraSerialNumber", dv.FilePath)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// parse output
	var data []map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return err
	}

	// Video Information
	if fileName, ok := data[0]["FileName"].(string); ok {
		dv.FileName = strings.Trim(fileName, ".bin")
	}

	if fileType, ok := data[0]["FileTypeExtension"].(string); ok {
		dv.FileType = fileType
	}

	if duration, ok := data[0]["Duration"].(string); ok {
		var durationFloat float64
		if strings.Contains(duration, "s") {
			duration = strings.Trim(duration, " s")
			durationFloat, err = strconv.ParseFloat(duration, 32)
			if err != nil {
				return err
			}
		} else {
			durationFloat, err = convertDurationToSeconds(duration)
			if err != nil {
				return err
			}
		}

		dv.Duration = tools.TruncateFloat(durationFloat)
	}

	if mediaUniqueID, ok := data[0]["MediaUniqueID"].(string); ok {
		dv.MediaUniqueID = mediaUniqueID
	}

	// In future, need to save Camera in independant table
	if cameraSerialNumber, ok := data[0]["CameraSerialNumber"].(string); ok {
		dv.CameraSerialNumber = cameraSerialNumber
	}

	if model, ok := data[0]["Model"].(string); ok {
		dv.CameraModel = model
	}

	return nil
}

func (dv *Video) FillRawMetadata(file []byte) ([]byte, error) {
	err := ioutil.WriteFile(dv.Name, file, 0644)
	if err != nil {
		return nil, err
	}

	// Use Exiftool CLI to get metadata, try to get Camera Model, Duration and FileName and FileType + MediaUniqueID + CameraSerialNumber
	cmd := exec.Command("exiftool", "-j", "-Model", "-Duration", "-FileName", "-FileTypeExtension", "-MediaUniqueID", "-CameraSerialNumber", dv.Name)
	output, err := cmd.Output()

	return output, err
	// fmt.Println(output)
	// if err != nil {
	// 	return err
	// }

	// // parse output
	// var data []map[string]interface{}
	// if err := json.Unmarshal(output, &data); err != nil {
	// 	return err
	// }

	// // Video Information
	// if fileName, ok := data[0]["FileName"].(string); ok {
	// 	dv.FileName = strings.Trim(fileName, ".MP4")
	// }

	// if fileType, ok := data[0]["FileTypeExtension"].(string); ok {
	// 	dv.FileType = fileType
	// }

	// if duration, ok := data[0]["Duration"].(string); ok {
	// 	var durationFloat float64
	// 	if strings.Contains(duration, "s") {
	// 		duration = strings.Trim(duration, " s")
	// 		durationFloat, err = strconv.ParseFloat(duration, 32)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		durationFloat, err = convertDurationToSeconds(duration)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	// 	dv.Duration = tools.TruncateFloat(durationFloat)
	// }

	// if mediaUniqueID, ok := data[0]["MediaUniqueID"].(string); ok {
	// 	dv.MediaUniqueID = mediaUniqueID
	// }

	// // In future, need to save Camera in independant table
	// if cameraSerialNumber, ok := data[0]["CameraSerialNumber"].(string); ok {
	// 	dv.CameraSerialNumber = cameraSerialNumber
	// }

	// if model, ok := data[0]["Model"].(string); ok {
	// 	dv.CameraModel = model
	// }

	// return nil
}

func (dv *Video) ExtractGpxDataFromBinFile(fileContent []byte) (Gpx, error) {
	//command : "gopro2gpx -i GOPR0001.bin -a 500 -f 2 -o GOPR0001.gpx"
	gopro2gpxPath := os.Getenv("GOPRO2GPX_PATH")
	gpxFilePath := os.Getenv("GPX_FILES_DEST_DIR") + dv.Name + ".gpx"
	tempFile, err := os.Create(os.Getenv("RAW_VIDEO_DEST_DIR") + dv.FileName)
	if err != nil {
		return Gpx{}, fmt.Errorf("error while creating temp file: %w", err)
	}
	defer tempFile.Close()

	_, err = tempFile.Write(fileContent)
	if err != nil {
		return Gpx{}, fmt.Errorf("error while writing to temp file: %w", err)
	}

	convertBinToGpx := exec.Command(gopro2gpxPath, "-i", os.Getenv("RAW_VIDEO_DEST_DIR")+dv.FileName, "-a", "500", "-o", gpxFilePath)
	err = convertBinToGpx.Run()
	if err != nil {
		err = fmt.Errorf("error while converting bin to gpx : %w", err)
		return Gpx{}, err
	}

	var gpx Gpx
	gpx.Name = dv.Name + ".gpx"
	gpx.Type = TypeFromGopro

	// Read GPX file
	gpxFileContent, err := os.ReadFile(gpxFilePath)
	if err != nil {
		err = fmt.Errorf("error while reading gpx file : %w", err)
		return Gpx{}, err
	}

	var gpxStruct GpxStruct
	err = xml.Unmarshal([]byte(gpxFileContent), &gpxStruct)
	if err != nil {
		err = fmt.Errorf("error while unmarshal gpx file : %w", err)
		return Gpx{}, err
	}

	// check if gpxStruct is empty
	if len(gpxStruct.Trk.Trkseg) == 0 ||
		len(gpxStruct.Trk.Trkseg[0].Trkpt) == 0 {
		return gpx, nil
	}
	// Get first Coord
	gpx.StartLat = gpxStruct.Trk.Trkseg[0].Trkpt[0].Lat
	gpx.StartLon = gpxStruct.Trk.Trkseg[0].Trkpt[0].Lon

	// Get last Coord
	gpx.EndLat = gpxStruct.Trk.Trkseg[0].Trkpt[len(gpxStruct.Trk.Trkseg[0].Trkpt)-1].Lat
	gpx.EndLon = gpxStruct.Trk.Trkseg[0].Trkpt[len(gpxStruct.Trk.Trkseg[0].Trkpt)-1].Lon

	return gpx, nil
}

func convertDurationToSeconds(duration string) (float64, error) {
	t, err := time.Parse("15:04:05", duration)
	if err != nil {
		return 0, err
	}

	hours := float64(t.Hour())
	minutes := float64(t.Minute())
	seconds := float64(t.Second())

	totalSeconds := hours*3600 + minutes*60 + seconds
	return totalSeconds, nil
}

func (dv *Video) TransformMp4FileToBinFile() error {
	rawFileDestDir := os.Getenv("RAW_VIDEO_DEST_DIR")
	cmd := exec.Command("ffmpeg", "-y", "-i", dv.FilePath, "-codec", "copy", "-map", "0:3", "-f", "rawvideo", rawFileDestDir+dv.FileName+".bin")

	return cmd.Run()
}
