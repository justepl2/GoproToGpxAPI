package response

import (
	"github.com/justepl2/gopro_to_gpx_api/domain"
)

type Video struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	FileName    string  `json:"fileName"`
	Duration    float64 `json:"duration"`
	CameraModel string  `json:"cameraModel"`
	Status      string  `json:"status"`
	Gpx         Gpx     `json:"gpx"`
}

func (r *Video) FromDomain(video domain.Video) {
	r.Id = video.ID.String()
	r.Name = video.Name
	r.FileName = video.FileName
	r.Duration = video.Duration
	r.CameraModel = video.CameraModel
	r.Status = string(video.Status)

	r.Gpx.FromDomain(video.Gpx)
}
