package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GpxType string

const (
	TypeRead     GpxType = "read"
	TypeGenerate GpxType = "generate"
)

type Gpx struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"column:name"`
	StartCoords string    `gorm:"column:start_coords"`
	EndCoords   string    `gorm:"column:end_coords"`
	S3Location  string    `gorm:"column:s3_location"`
	Type        GpxType   `gorm:"column:type"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
