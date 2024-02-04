package videos

import (
	"encoding/json"
	"net/http"

	"github.com/justepl2/gopro_to_gpx_api/application"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/response"
)

func List(w http.ResponseWriter, r *http.Request) {
	videos, err := application.ListVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	videosResponse := make([]response.Video, len(videos))
	for i, video := range videos {
		videosResponse[i].FromDomain(video)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videosResponse)
}
