package main

import (
	"fmt"

	"github.com/justepl2/gopro_to_gpx_api/interfaces"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting server...")
	godotenv.Load(".env.local")
	interfaces.Run(8080)
}
