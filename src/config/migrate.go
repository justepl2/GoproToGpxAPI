package config

import (
	"log"

	"github.com/justepl2/gopro_to_gpx_api/domain"
)

func DBMigrate() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}

	db.AutoMigrate(&domain.Video{}, &domain.Gpx{})
	if db.Error != nil {
		log.Println("Error migrating database")
		return db.Error
	}

	log.Println("Migration has been processed")

	return nil
}
