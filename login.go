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
		Lifetime int    `json:"expires_in_seconds"`
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

	newToken, err := cfg.createNewToken(loginUser, params.Lifetime)
	if err != nil {
		log.Printf("unable to create JWT: %v\n", err)
		w.WriteHeader(500)
		return
	}

	loginUserResp := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		ID:    loginUser.ID,
		Email: loginUser.Email,
		Token: newToken,
	}

	respondWithJSON(w, 200, loginUserResp)
}
