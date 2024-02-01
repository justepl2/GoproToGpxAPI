package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"gorm.io/gorm"
)

type Video struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name       string    `gorm:"column:name"`
	Duration   float64   `gorm:"column:duration"`
	S3Location string    `gorm:"column:s3_location"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (dv *Video) FromRequest(rv request.CreateVideo) {
	dv.Name = rv.Name
	dv.Duration = rv.Duration
	dv.S3Location = "bucket/default/location"
}
