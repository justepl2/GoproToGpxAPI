package application

import (
	"github.com/justepl2/gopro_to_gpx_api/config"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure/db"
)

func AddGpx(gpx *domain.Gpx) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := db.NewGpxRepository(conn)

	return repo.Save(gpx)
}

func GetGpxById(id string) (*domain.Gpx, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := db.NewGpxRepository(conn)
	return repo.FindById(id)
}
