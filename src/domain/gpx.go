package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GpxType string

const (
	TypeFromGopro GpxType = "FromGopro"
	TypeLinker    GpxType = "Linker"
)

type Gpx struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	// VideoID    string    `gorm:"column:video_id"`
	Name       string  `gorm:"column:name"`
	StartLat   string  `gorm:"column:start_lat"`
	StartLon   string  `gorm:"column:start_lon"`
	EndLat     string  `gorm:"column:end_lat"`
	EndLon     string  `gorm:"column:end_lon"`
	S3Location string  `gorm:"column:s3_location"`
	Type       GpxType `gorm:"column:type"`
	Status     Status  `gorm:"column:status"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Gpx) TableName() string {
	return "gpx"
}
