package http

import (
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, "failed to encode json", http.StatusInternalServerError)
		}
	}
}
