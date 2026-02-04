package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {

	type envelope struct {
		Error string `json:"error"`
	}

	data := envelope{
		Error: message,
	}

	return writeJSON(w, status, data)

}

func readJSON(r *http.Request, dst any) error {

	maxBytes := 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(nil, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)

	// Disallow unknown fields
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}
