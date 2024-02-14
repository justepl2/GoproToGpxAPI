package videos

import (
	"encoding/json"
	"fmt"
	"io"
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

	fmt.Println("endpoint POST /videos called")
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

	// Create Raw Video file
	err = video.TransformMp4FileToBinFile()
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot transform video to bin file, err : "+err.Error())
		return
	}

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
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot upload file to S3, err : "+err.Error())
		return
	}

	// Extract GPX Data from Raw Video file
	// gpx, err := video.ExtractGpxDataFromBinFile()
	// if err != nil {
	// tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot extract GPX data from bin file, err : "+err.Error())
	// return
	// }

	// push GPX on S3
	// gpxFileContent, err := os.ReadFile(os.Getenv("GPX_FILES_DEST_DIR") + video.FileName + ".gpx")
	// if err != nil {
	// 	tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot read gpx file, err : "+err.Error())
	// 	return
	// }

	// err = s3.UploadFiles(video.CameraSerialNumber+"/"+video.FileName+".gpx", gpxFileContent)
	// if err != nil {
	// 	tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot upload gpx file to S3, err : "+err.Error())
	// 	return
	// }

	// // Create GPX on DB
	// gpx.S3Location = video.CameraSerialNumber + "/" + video.FileName + ".gpx"
	// gpx.Status = domain.StatusDone
	// err = application.AddGpx(&gpx)
	// if err != nil {
	// 	tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot add gpx to Database, err : "+err.Error())
	// 	return
	// }

	// // Create Video on DB
	// video.S3Location = video.CameraSerialNumber + "/" + video.FileName + ".bin"
	// video.Gpx = gpx
	// video.Status = domain.StatusDone
	// err = application.AddVideo(&video)
	// if err != nil {
	// 	tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot add video to Database, err : "+err.Error())
	// 	return
	// }

	// // Delete Raw Video file
	// tools.DeleteTempFiles(video.FileName)
	// tools.FormatResponseBody(w, http.StatusCreated, video.ID.String())
}

func CreateFromRaw(w http.ResponseWriter, r *http.Request) {
	var video domain.Video

	// Parse the multipart form in the request
	err := r.ParseMultipartForm(10 << 20) // Max memory 10MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the file from form data
	file, fileHeader, err := r.FormFile("file") // "file" is the key of the form-data to be uploaded
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "Cannot retrieve file from form-data, err : "+err.Error())
		return
	}
	defer file.Close()

	// Read the file into a byte slice (byte[])
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot read file, err : "+err.Error())
		return
	}

	// Create a RawFile
	rawFile := request.RawFile{
		Name: fileHeader.Filename,
		File: fileBytes,
	}

	rawFile.Validate()

	video.FromRawRequest(rawFile)
	// output, err := video.FillRawMetadata(fileBytes)
	//if err != nil {
	//	tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot get video metadata, err : "+err.Error())
	//	return
	//}

	gpx, err := video.ExtractGpxDataFromBinFile(rawFile.File)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot extract GPX data from bin file, err : "+err.Error())
		return
	}

	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot read gpx file, err : "+err.Error())
		return
	}

	gpxFileContent, err := os.ReadFile(os.Getenv("GPX_FILES_DEST_DIR") + video.Name + ".gpx")
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot read gpx file, err : "+err.Error())
		return
	}

	s3 := infrastructure.NewS3FileStorage()
	err = s3.UploadFiles(video.CameraSerialNumber+"/"+video.FileName+".gpx", gpxFileContent)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot upload gpx file to S3, err : "+err.Error())
		return
	}

	// Create GPX on DB
	gpx.S3Location = video.CameraSerialNumber + "/" + video.FileName + ".gpx"
	gpx.Status = domain.StatusDone
	err = application.AddGpx(&gpx)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot add gpx to Database, err : "+err.Error())
		return
	}

	// Create Video on DB
	video.S3Location = video.CameraSerialNumber + "/" + video.FileName + ".bin"
	video.Gpx = gpx
	video.Status = domain.StatusDone
	err = application.AddVideo(&video)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot add video to Database, err : "+err.Error())
		return
	}

	tools.FormatResponseBody(w, http.StatusCreated, "File received")
}
