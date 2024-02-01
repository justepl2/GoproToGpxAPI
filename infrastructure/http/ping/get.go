package ping

import (
	"encoding/json"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "pong"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
