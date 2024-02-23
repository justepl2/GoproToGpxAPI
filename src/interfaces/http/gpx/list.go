package gpx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/response"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

// List godoc
// @Summary List all GPX
// @Description List all GPX
// @Tags gpx
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} response.Gpx "OK"
// @Failure 500 {object} response.Error "Internal server error"
// @Router /gpx [get]
func List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint GET /gpx called")

	gpxs, err := application.ListGpx()
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	gpxsResponse := make([]response.Gpx, len(gpxs))
	for i, gpx := range gpxs {
		gpxsResponse[i].FromDomain(gpx)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gpxsResponse)
}
