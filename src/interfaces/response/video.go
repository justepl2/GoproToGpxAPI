package response

import (
	"github.com/justepl2/gopro_to_gpx_api/domain"
)

type Video struct {
	Id          string  `json:"id" example:"5f5e3e4e-3e4e-5f5e-3e4e-5f5e3e4e3e4e"`
	Name        string  `json:"name" example:"video_1"`
	FileName    string  `json:"fileName" example:"video_1.bin"`
	Duration    float64 `json:"duration" example:"10.0"`
	CameraModel string  `json:"cameraModel" example:"GoPro Hero 8"`
	Status      string  `json:"status" example:"FromGopro"`
	Gpx         Gpx     `json:"gpx"`
}

func (r *Video) FromDomain(video domain.Video) {
	r.Id = video.ID.String()
	r.Name = video.Name
	r.FileName = video.FileName
	r.Duration = video.Duration
	r.CameraModel = video.CameraModel
	r.Status = string(video.Status)

	r.Gpx.FromDomain(video.Gpx, "")
}
