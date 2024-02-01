package interfaces

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure/http/ping"
)

// Run start server
func Run(port int) error {
	log.Printf("Server running at http://localhost:%d/", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), NewRouter())
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", ping.Get).Methods("GET")
	// router.HandleFunc("/gpx", handler.HandleRequest).Methods("GET")

	return router
}
