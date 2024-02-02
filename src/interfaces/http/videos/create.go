package videos

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

func Create(w http.ResponseWriter, r *http.Request) {
	var requestVideo request.CreateVideo
	var video domain.Video

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestVideo)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	err = requestVideo.Validate()
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, err.Error())
		return
	}

	video.FromRequest(requestVideo)

	// get Video Metadata
	err = video.FillVideoMetadata()
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot get video metadata, err : "+err.Error())
		return
	}

	// Create Video on DB
	video.Status = domain.StatusPending
	err = application.AddVideo(&video)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot add video to Database, err : "+err.Error())
		return
	}

	// Create Raw Video file
	err = video.TransformMp4FileToBinFile()
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot transform video to bin file, err : "+err.Error())
		return
	}
	defer tools.DeleteTempFiles(video.FileName)

	// Push on S3
	s3 := infrastructure.NewS3FileStorage()
	filePath := os.Getenv("RAW_VIDEO_DEST_DIR") + video.FileName
	fileContent, err := os.ReadFile(filePath + ".bin")
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot read file, err : "+err.Error())
		return
	}

	err = s3.UploadFiles(video.CameraSerialNumber+"/"+video.FileName+".bin", fileContent)
	if err != nil {
		// if error, update DB with status Error
		handleError(w, &video, err, "Cannot upload file to S3")
		return
	}

	// Extract GPX Data from Raw Video file
	gpx, err := video.ExtractGpxDataFromBinFile()
	if err != nil {
		handleError(w, &video, err, "Cannot extract GPX data from bin file")
		return
	}

	// push GPX on S3
	gpxFileContent, err := json.Marshal(os.Getenv("GPX_FILES_DEST_DIR") + video.FileName + ".gpx")
	if err != nil {
		handleError(w, &video, err, "Cannot read gpx file")
		return
	}

	err = s3.UploadFiles(video.CameraSerialNumber+"/"+video.FileName+".gpx", gpxFileContent)
	if err != nil {
		handleError(w, &video, err, "Cannot upload gpx file to S3")
		return
	}

	// Create GPX on DB
	gpx.S3Location = video.CameraSerialNumber + "/" + video.FileName + ".gpx"
	gpx.Status = domain.StatusDone
	err = application.AddGpx(&gpx)
	if err != nil {
		handleError(w, &video, err, "Cannot add gpx to Database")
		return
	}

	// Update Video on DB
	video.S3Location = video.CameraSerialNumber + "/" + video.FileName + ".bin"
	video.Status = domain.StatusDone
	err = application.UpdateVideo(&video)
	if err != nil {
		handleError(w, &video, err, "Cannot update video status")
		return
	}

	// Delete Raw Video file
	tools.FormatResponseBody(w, http.StatusCreated, video.ID.String())
}

func handleError(w http.ResponseWriter, video *domain.Video, err error, message string) {
	video.Status = domain.StatusError
	updateErr := application.UpdateVideo(video)
	if updateErr != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot set video status to 'Error', err : "+updateErr.Error())
		return
	}
	tools.FormatResponseBody(w, http.StatusInternalServerError, message+", err : "+err.Error())
}