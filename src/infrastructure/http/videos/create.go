package videos

import (
	"encoding/json"
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/domain"
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

	video.FromRequest(requestVideo)
	err = application.AddVideo(&video)
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, "Cannot add video to Database, err : "+err.Error())
		return
	}

	tools.FormatResponseBody(w, http.StatusCreated, video.ID.String())
}
