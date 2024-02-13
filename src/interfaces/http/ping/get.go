package ping

import (
	"fmt"
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/tools"
)

func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint GET /ping called")
	tools.FormatResponseBody(w, http.StatusOK, "pong")
}
