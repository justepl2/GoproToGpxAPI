package interfaces

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure/http/ping"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure/http/videos"
)

// Run start server
func Run(port int) error {
	log.Printf("Server running at http://localhost:%d/", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), NewRouter())
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	// HealthCheck endpoint
	router.HandleFunc("/ping", ping.Get).Methods("GET")

	// Videos endpoints
	router.HandleFunc("/video", videos.Create).Methods("POST")
	router.HandleFunc("/video", videos.Create).Methods("GET")

	// GPX endpoints
	return router
}
