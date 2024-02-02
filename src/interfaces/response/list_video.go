package response

import "github.com/justepl2/gopro_to_gpx_api/domain"

type ListVideo struct {
	Name        string  `json:"name"`
	FileName    string  `json:"fileName"`
	Duration    float64 `json:"duration"`
	CameraModel string  `json:"cameraModel"`
	Status      string  `json:"status"`
}

func (lv *ListVideo) FromDomain(video domain.Video) {
	lv.Name = video.Name
	lv.FileName = video.FileName
	lv.Duration = video.Duration
	lv.CameraModel = video.CameraModel
	lv.Status = string(video.Status)
}
