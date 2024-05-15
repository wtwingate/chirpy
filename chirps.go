package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// respond with an Array of all chirps, sorted by ID
}

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %v\n", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		cleanBody := filterProfanity(params.Body)
		chirp, err := cfg.db.NewChirp(cleanBody)
		if err != nil {
			log.Printf("error posting chirp: %v\n", err)
			w.WriteHeader(500)
			return
		}
		respondWithJSON(w, 201, chirp)
	}
}

func filterProfanity(msg string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range strings.Fields(msg) {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				cleaned := strings.Split(msg, word)
				msg = strings.Join(cleaned, "****")
			}
		}
	}
	return msg
}
