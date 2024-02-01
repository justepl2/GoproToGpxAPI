package db

import (
	"github.com/jinzhu/gorm"
	"github.com/justepl2/gopro_to_gpx_api/domain"
)

// NewsRepositoryImpl Implements repository.NewsRepository
type NewsRepositoryImpl struct {
	Conn *gorm.DB
}

func NewVideoRepository(conn *gorm.DB) *NewsRepositoryImpl {
	return &NewsRepositoryImpl{Conn: conn}
}

func (r *NewsRepositoryImpl) Save(video *domain.Video) error {
	if err := r.Conn.Save(video).Error; err != nil {
		return err
	}

	return nil
}
