package main

import (
	"encoding/json"
	"net/http"

	"github.com/wtwingate/chirpy/internal/auth"
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
		respondWithError(w, http.StatusInternalServerError, "unable to decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to find user")
		return
	}

	err = auth.CheckHashPassword(user.Hash, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	authToken, err := auth.CreateNewAuthToken(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create auth token")
		return
	}

	refreshToken, err := cfg.DB.CreateNewRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create refresh token")
	}

	loginUserResp := struct {
		ID      int    `json:"id"`
		Email   string `json:"email"`
		Token   string `json:"token"`
		Refresh string `json:"refresh_token"`
	}{
		ID:      user.ID,
		Email:   user.Email,
		Token:   authToken,
		Refresh: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, loginUserResp)
}
