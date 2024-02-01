package main

import (
	"fmt"

	"github.com/justepl2/gopro_to_gpx_api/config"
	"github.com/justepl2/gopro_to_gpx_api/interfaces"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting server...")
	godotenv.Load(".env.local")
	err := config.DBMigrate()
	if err != nil {
		panic(err)
	}

	interfaces.Run(8080)
}
