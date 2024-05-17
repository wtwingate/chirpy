package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %v\n", err)
		w.WriteHeader(500)
		return
	}

	if _, err = cfg.db.GetUserByEmail(params.Email); err == nil {
		log.Println("user already exists with that email")
		w.WriteHeader(401)
		return
	}

	newUser, err := cfg.db.NewUser(params.Email, params.Password)
	if err != nil {
		log.Printf("error creating new user: %v\n", err)
		w.WriteHeader(500)
		return
	}

	newUserResp := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    newUser.ID,
		Email: newUser.Email,
	}

	respondWithJSON(w, 201, newUserResp)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string
		Password string
	}
}
