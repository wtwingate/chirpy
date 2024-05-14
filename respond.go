package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	log.Printf("Responding with status code %v error: %v", code, msg)
	errorResponse := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}
	respondWithJSON(w, code, errorResponse)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error: unable to marshal JSON")
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
