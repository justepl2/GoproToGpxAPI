package interfaces

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure/middleware"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/http/gpx"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/http/ping"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/http/users"
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

	// Users endpoints
	router.HandleFunc("/users/signup", users.Signup).Methods("POST")
	router.HandleFunc("/users/login", users.Login).Methods("POST")
	router.HandleFunc("/users/me", users.Me).Methods("GET")
	router.Handle("/users/logout", middleware.UserAuthenticationMiddleware(http.HandlerFunc(users.Logout))).Methods("POST")
	router.Handle("/users/refresh", middleware.UserAuthenticationMiddleware(http.HandlerFunc(users.Refresh))).Methods("POST")
	router.HandleFunc("/users/forgot-password", users.ForgotPassword).Methods("POST")

	// Videos endpoints
	router.Handle("/videos", middleware.UserAuthenticationMiddleware(http.HandlerFunc(videos.List))).Methods("GET")
	router.Handle("/videos/link", middleware.UserAuthenticationMiddleware(http.HandlerFunc(videos.Link))).Methods("POST")
	router.Handle("/videos/raw", middleware.UserAuthenticationMiddleware(http.HandlerFunc(videos.CreateFromRaw))).Methods("POST")

	// GPX endpoints
	router.Handle("/gpx", middleware.UserAuthenticationMiddleware(http.HandlerFunc(gpx.List))).Methods("GET")
	router.Handle("/gpx/{id}", middleware.UserAuthenticationMiddleware(http.HandlerFunc(gpx.GetById))).Methods("GET")

	// Apply the CORS middleware to our top-level router, with the defaults.
	router.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
	))

	return router
}
