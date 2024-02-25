package tools

import (
	"encoding/json"
	"net/http"
)

func FormatEmptyResponseBody(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{}`))
}

func FormatResponseBody(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]string{"message": message}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func FormatUuidResponseBody(w http.ResponseWriter, statusCode int, uuid string) {
	response := map[string]string{"uuid": uuid}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func FormatStrResponseBody(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
