package db

import (
	"gorm.io/gorm"

	"github.com/justepl2/gopro_to_gpx_api/domain"
)

type VideoRepositoryImpl struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryImpl {
	return &VideoRepositoryImpl{db: db}
}

func (r *VideoRepositoryImpl) Save(video *domain.Video) error {
	return r.db.Save(video).Error
}

func (r *VideoRepositoryImpl) Update(video *domain.Video) error {
	return r.db.Model(video).Updates(video).Error
}

func (r *VideoRepositoryImpl) FindAll() ([]domain.Video, error) {
	var videos []domain.Video
	err := r.db.Joins("Gpx").Find(&videos).Error
	return videos, err
}
