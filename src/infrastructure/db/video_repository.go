package db

import (
	"github.com/jinzhu/gorm"
	"github.com/justepl2/gopro_to_gpx_api/domain"
)

type VideoRepositoryImpl struct {
	Conn *gorm.DB
}

func NewVideoRepository(conn *gorm.DB) *VideoRepositoryImpl {
	return &VideoRepositoryImpl{Conn: conn}
}

func (r *VideoRepositoryImpl) Save(video *domain.Video) error {
	return r.Conn.Save(video).Error
}

func (r *VideoRepositoryImpl) Update(video *domain.Video) error {
	return r.Conn.Model(video).Updates(video).Error
}

func (r *VideoRepositoryImpl) FindAll() ([]domain.Video, error) {
	var videos []domain.Video
	err := r.Conn.Find(&videos).Error
	return videos, err
}
