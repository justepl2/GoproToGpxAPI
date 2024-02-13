package response

import (
	"github.com/google/uuid"
	"github.com/justepl2/gopro_to_gpx_api/domain"
)

type Gpx struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	StartCoords map[string]string `json:"startCoords"`
	EndCoords   map[string]string `json:"endCoords"`
	S3Location  string            `json:"s3Location"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
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
