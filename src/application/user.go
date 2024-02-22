package application

import (
	"github.com/justepl2/gopro_to_gpx_api/config"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/infrastructure/db"
)

func AddUser(user *domain.User) error {
	conn, err := config.ConnectDB()
	if err != nil {
		return err
	}

	repo := db.NewUserRepository(conn)

	return repo.Save(user)
}

func GetUserByEmail(email string) (domain.User, error) {
	conn, err := config.ConnectDB()
	if err != nil {
		return domain.User{}, err
	}

	repo := db.NewUserRepository(conn)

	return repo.GetByEmail(email)
}
