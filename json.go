package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with %v error: %v", code, msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(payload); err != nil {
		log.Printf("Failed to encode JSON response: %v, payload: %v", err, payload)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
}
