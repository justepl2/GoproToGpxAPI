package tools

import (
	"encoding/json"
	"net/http"
)

func FormatResponseBody(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]string{"message": message}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
