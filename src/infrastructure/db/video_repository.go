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
	return r.Conn.Save(video).Error
}

func (r *NewsRepositoryImpl) Update(video *domain.Video) error {
	return r.Conn.Model(video).Updates(video).Error
}
