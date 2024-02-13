package interfaces

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/http/gpx"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/http/ping"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/http/videos"
)

// Run start server
func Run(port int) error {
	log.Printf("Server running at http://localhost:%d/", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), NewRouter())
	if err != nil {
		return err
	}

	return nil
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	// HealthCheck endpoint
	router.HandleFunc("/ping", ping.Get).Methods("GET")

	// Videos endpoints
	router.HandleFunc("/video", videos.Create).Methods("POST")
	router.HandleFunc("/video", videos.List).Methods("GET")
	router.HandleFunc("/video/link", videos.Link).Methods("POST")

	// GPX endpoints
	router.HandleFunc("/gpx", gpx.List).Methods("GET")
	router.HandleFunc("/gpx/{id}", gpx.GetById).Methods("GET")

	return router
}
