package db

import (
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"gorm.io/gorm"
)

type GpxRepositoryImpl struct {
	Conn *gorm.DB
}

func NewGpxRepository(conn *gorm.DB) *GpxRepositoryImpl {
	return &GpxRepositoryImpl{Conn: conn}
}

func (r *GpxRepositoryImpl) Save(gpx *domain.Gpx) error {
	return r.Conn.Create(&gpx).Error
}

func (r *GpxRepositoryImpl) FindAll() ([]domain.Gpx, error) {
	var gpx []domain.Gpx
	err := r.Conn.Find(&gpx).Error
	return gpx, err
}

func (r *GpxRepositoryImpl) FindById(id string) (*domain.Gpx, error) {
	var gpx domain.Gpx
	err := r.Conn.Where("id = ?", id).First(&gpx).Error
	return &gpx, err
}
