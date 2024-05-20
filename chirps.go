package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/wtwingate/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirpByID(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid chirp ID")
		return
	}

	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		respondWithError(w, http.StatusUnauthorized, "missing authorization token")
		return
	}

	userID, err := auth.CheckAuthToken(cfg.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	err = cfg.DB.DeleteChirp(chirpID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode parameters")
		return
	}

	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		respondWithError(w, http.StatusUnauthorized, "missing authorization token")
		return
	}

	userID, err := auth.CheckAuthToken(cfg.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "chirp is too long")
	} else {
		cleanBody := filterProfanity(params.Body)
		chirp, err := cfg.DB.CreateChirp(userID, cleanBody)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, chirp)
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
