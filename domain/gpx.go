package domain

import "github.com/google/uuid"

type GpxType string

const (
	TypeRead     GpxType = "read"
	TypeGenerate GpxType = "generate"
)

type Gpx struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string    `gorm:"column:name"`
	StartCoords string    `gorm:"column:start_coords"`
	EndCoords   string    `gorm:"column:end_coords"`
	S3Location  string    `gorm:"column:s3_location"`
	Type        GpxType   `gorm:"column:type"`
}
