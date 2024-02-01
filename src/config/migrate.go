package config

import (
	"log"

	"github.com/justepl2/gopro_to_gpx_api/domain"
)

func DBMigrate() error {
	conn, err := ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.AutoMigrate(&domain.Video{}, &domain.Gpx{})
	if conn.Error != nil {
		log.Println("Error migrating database")
		return conn.Error
	}

	log.Println("Migration has been processed")

	return nil
}
