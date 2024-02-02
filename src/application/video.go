package application

import (
	"github.com/justepl2/gopro_to_gpx_api/config"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure/db"
)

func AddVideo(video *domain.Video) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := db.NewVideoRepository(conn)

	return repo.Save(video)
}

func UpdateVideo(video *domain.Video) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	repo := db.NewVideoRepository(conn)

	return repo.Update(video)
}

func ListVideos() ([]domain.Video, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := db.NewVideoRepository(conn)

	return repo.FindAll()
}
