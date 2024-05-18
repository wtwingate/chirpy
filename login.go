package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wtwingate/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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

	refreshToken, err := cfg.DB.CreateNewRefreshToken(user.ID)
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

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refToken := r.Header.Get("Authorization")
	if len(refToken) == 0 {
		respondWithError(w, http.StatusUnauthorized, "missing authorization token")
		return
	}

	refToken = strings.TrimPrefix(refToken, "Bearer ")

	userID, err := cfg.DB.CheckRefreshToken(refToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
	}

	authToken, err := auth.CreateNewAuthToken(userID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create auth token")
		return
	}

	refreshResp := struct {
		Token string `json:"token"`
	}{
		Token: authToken,
	}

	respondWithJSON(w, 200, refreshResp)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refToken := r.Header.Get("Authorization")
	if len(refToken) == 0 {
		respondWithError(w, http.StatusUnauthorized, "missing authorization token")
		return
	}

	refToken = strings.TrimPrefix(refToken, "Bearer ")

	err := cfg.DB.RevokeRefreshToken(refToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to revoke refresh token")
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
