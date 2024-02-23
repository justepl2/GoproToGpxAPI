package users

import (
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/tools"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	tools.FormatResponseBody(w, http.StatusOK, "Logout endpoint")
}
