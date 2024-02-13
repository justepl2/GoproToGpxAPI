package gpx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

func GetById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint GET /gpx/{id} called")

	vars := mux.Vars(r)
	gpx, err := application.GetGpxById(vars["id"])
	if err != nil {
		tools.FormatResponseBody(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gpx)
}
