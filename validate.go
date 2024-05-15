package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidation(w http.ResponseWriter, r *http.Request) {
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
		responseBody := filterProfanity(params.Body)
		respondWithJSON(w, 200, struct {
			CleanedBody string `json:"cleaned_body"`
		}{
			CleanedBody: responseBody,
		})
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
