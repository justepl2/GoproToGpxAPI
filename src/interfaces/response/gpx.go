package response

import (
	"github.com/google/uuid"
	"github.com/justepl2/gopro_to_gpx_api/domain"
)

type Gpx struct {
	ID          uuid.UUID         `json:"id" example:"5f5e3e4e-3e4e-5f5e-3e4e-5f5e3e4e3e4e"`
	Name        string            `json:"name" example:"video_1"`
	StartCoords map[string]string `json:"startCoords" example:"lat: 0.0, lon: 0.0"`
	EndCoords   map[string]string `json:"endCoords" example:"lat: 0.0, lon: 0.0"`
	S3Location  string            `json:"s3Location" example:"s3://bucket/folder/file.gpx"`
	Type        string            `json:"type" example:"road"`
	Status      string            `json:"status" example:"FromGopro"`
}

func (r *Gpx) FromDomain(dgpx domain.Gpx) {
	r.ID = dgpx.ID
	r.Name = dgpx.Name
	r.StartCoords = map[string]string{"lat": dgpx.StartLat, "lon": dgpx.StartLon}
	r.EndCoords = map[string]string{"lat": dgpx.EndLat, "lon": dgpx.EndLon}
	r.S3Location = dgpx.S3Location
	r.Type = string(dgpx.Type)
	r.Status = string(dgpx.Status)
}
