package gpx

import (
	"encoding/json"
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

func List(w http.ResponseWriter, r *http.Request) {
	gpxs, err := application.ListGpx()
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gpxs)
}