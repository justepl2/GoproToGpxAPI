package videos

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

func Link(w http.ResponseWriter, r *http.Request) {
	// Get from body the video id to link
	// should have at least 2 videos
	var requestLinkVideos request.LinkVideos
	var gpxStruct domain.GpxStruct
	var responseGpxs []domain.Gpx

	fmt.Println("endpoint POST /videos/link called")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestLinkVideos)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	err = requestLinkVideos.Validate()
	if err != nil {
		tools.FormatResponseBody(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get Videos.gpx from DB
	videos, err := application.GetVideosByIds(requestLinkVideos.VideoIds)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	// for each gpx
	for i := 0; i < len(videos)-1; i++ {
		if videos[i].Gpx == (domain.Gpx{}) || videos[i+1].Gpx == (domain.Gpx{}) {
			tools.FormatResponseBody(w, http.StatusBadRequest, "given video doesn't contain gpx data, could be a Gopro error")
			return
		}

		// get get gpx[i].EndCoords
		// and get gpx[i+1].StartCoords
		video1LastCoord := [2]string{videos[i].Gpx.EndLat, videos[i].Gpx.EndLon}
		video2FirstCoord := [2]string{videos[i+1].Gpx.StartLat, videos[i+1].Gpx.StartLon}

		// call OpenRouteService API to get the route between the 2 points
		ors := infrastructure.NewOpenRouteService(string(requestLinkVideos.Terrain))
		linkRoute, err := ors.GetRoute(video1LastCoord, video2FirstCoord)
		if err != nil {
			tools.FormatResponseBody(w, http.StatusFailedDependency, err.Error())
			return
		}

		gpxStruct.ConvertOpenRouteServiceGpxIntoTrkGpx(linkRoute)

		content, err := xml.MarshalIndent(gpxStruct, "", " ")
		if err != nil {
			tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Push on S3
		s3 := infrastructure.NewS3FileStorage()
		err = s3.UploadFiles(videos[i].CameraSerialNumber+"/"+videos[i].FileName+"_to_"+videos[i+1].FileName+".gpx", content)
		if err != nil {
			tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Create Gpx in DB
		gpx := createGpxFromVideos(videos[i], videos[i+1])
		err = application.AddGpx(&gpx)
		if err != nil {
			tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
			return
		}

		responseGpxs = append(responseGpxs, gpx)
	}

	responseJson, err := json.Marshal(responseGpxs)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

func createGpxFromVideos(video1 domain.Video, video2 domain.Video) domain.Gpx {
	return domain.Gpx{
		Name:       video1.FileName + "_to_" + video2.FileName + ".gpx",
		StartLat:   video1.Gpx.StartLat,
		StartLon:   video1.Gpx.StartLon,
		EndLat:     video2.Gpx.EndLat,
		EndLon:     video2.Gpx.EndLon,
		S3Location: video1.CameraSerialNumber + "/" + video1.FileName + "_to_" + video2.FileName + ".gpx",
		Type:       domain.TypeLinker,
		Status:     domain.StatusDone,
	}
}
