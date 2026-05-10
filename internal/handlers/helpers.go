package handlers

import (
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respond(w, status, map[string]string{"error": msg})
}

func decode(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
