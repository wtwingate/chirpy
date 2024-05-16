package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
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

	loginUser, err := cfg.db.LoginUser(params.Email, params.Password)
	if err != nil {
		log.Printf("login error: %v\n", err)
		w.WriteHeader(401)
		return
	}

	loginUserResp := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    loginUser.ID,
		Email: loginUser.Email,
	}

	respondWithJSON(w, 200, loginUserResp)
}
