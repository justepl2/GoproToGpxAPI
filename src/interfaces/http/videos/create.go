package videos

import (
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

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

	userIdStr := r.Context().Value("userId").(string)
	rawFile.UserId, err = uuid.Parse(userIdStr)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	rawFile.Validate()

	video.FromRawRequest(rawFile)

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
	err = s3.UploadFiles(video.CameraSerialNumber+"/"+video.Name+".gpx", gpxFileContent)
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
	video.S3Location = video.UserId.String() + "/" + video.FileName + ".bin"
	video.Gpx = gpx
	video.Status = domain.StatusDone
	err = application.AddVideo(&video)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot add video to Database, err : "+err.Error())
		return
	}

	// Delete Raw Video file
	err = tools.DeleteTempFiles(video.Name)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot delete temp files, err : "+err.Error())
	}

	tools.FormatStrResponseBody(w, http.StatusCreated, video.ID.String())
}
