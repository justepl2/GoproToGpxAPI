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
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Documentation endpoint
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Users endpoints
	router.HandleFunc("/users/signup", users.Signup).Methods("POST")
	router.HandleFunc("/users/login", users.Login).Methods("POST")
	router.HandleFunc("/users/validateMail", users.ValidateMail).Methods("POST")
	router.HandleFunc("/users/forgot_password", users.ForgotPassword).Methods("POST")
	router.HandleFunc("/users/reset_password", users.ResetPassword).Methods("POST")
	router.HandleFunc("/users/refresh", users.Refresh).Methods("POST")
	router.Handle("/users/logout", middleware.LogoutUserMiddleware(http.HandlerFunc(users.Logout))).Methods("POST")

	// Videos endpoints
	router.Handle("/videos", middleware.ValidateTokenMiddleware(http.HandlerFunc(videos.List))).Methods("GET")
	router.Handle("/videos/link", middleware.ValidateTokenMiddleware(http.HandlerFunc(videos.Link))).Methods("POST")
	router.Handle("/videos/raw", middleware.ValidateTokenMiddleware(http.HandlerFunc(videos.CreateFromRaw))).Methods("POST")

	// GPX endpoints
	router.Handle("/gpx", middleware.ValidateTokenMiddleware(http.HandlerFunc(gpx.List))).Methods("GET")
	router.Handle("/gpx/{id}", middleware.ValidateTokenMiddleware(http.HandlerFunc(gpx.GetById))).Methods("GET")

	// Apply the CORS middleware to our top-level router, with the defaults.
	router.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
	))

	return router
}
