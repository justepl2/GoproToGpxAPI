package main

import (
	"fmt"
	"os"

	"github.com/justepl2/gopro_to_gpx_api/config"
	"github.com/justepl2/gopro_to_gpx_api/interfaces"

	"github.com/joho/godotenv"
)

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
