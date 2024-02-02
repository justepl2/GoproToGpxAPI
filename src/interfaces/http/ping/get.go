package ping

import (
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/tools"
)

func Get(w http.ResponseWriter, r *http.Request) {
	tools.FormatResponseBody(w, http.StatusOK, "pong")
}
