package domain

import "github.com/google/uuid"

type Video struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name       string    `gorm:"column:name"`
	Duration   float64   `gorm:"column:duration"`
	S3Location string    `gorm:"column:s3_location"`
}
