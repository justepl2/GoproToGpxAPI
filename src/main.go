package main

import (
	"fmt"
	"os"

	"github.com/justepl2/gopro_to_gpx_api/config"
	"github.com/justepl2/gopro_to_gpx_api/interfaces"

	"github.com/joho/godotenv"
	_ "github.com/justepl2/gopro_to_gpx_api/docs"
)

// @title Gopro GPX Extractor API
// @version 0.1
// @description This API extract GPX from Raw vidéo file and manage it.
// @termsOfService http://swagger.io/terms/

// @contact.name Mekadev
// @contact.url https://mekanull.com
// @contact.email contact@mekanull.com

// @host localhost:8080

// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	fmt.Println(os.Getenv("AWS_ACCESS_KEY_ID"))
	fmt.Println(os.Getenv("AWS_SECRET_ACCESS_KEY"))

	fmt.Println("Starting server...")
	godotenv.Load(".env.local")
	err := config.DBMigrate()
	if err != nil {
		panic(err)
	}

	interfaces.Run(8080)
}
