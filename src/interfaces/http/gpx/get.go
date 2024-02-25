package gpx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/response"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

// Get godoc
// @Summary GPX by ID
// @Description GPX by ID
// @Tags gpx
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "GPX ID"
// @Success 200 {object} response.Gpx "OK"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /gpx/{id} [get]
func GetById(w http.ResponseWriter, r *http.Request) {
	var gpxResponse response.Gpx

	fmt.Println("endpoint GET /gpx/{id} called")

	vars := mux.Vars(r)
	gpx, err := application.GetGpxById(vars["id"])
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(gpx.S3Location),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	gpxResponse.FromDomain(*gpx, urlStr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gpxResponse)
}
