package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wtwingate/chirpy/internal/auth"
	"github.com/wtwingate/chirpy/internal/database"
)

type request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Red      bool   `json:"is_chirpy_red"`
}

type response struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Red   bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := request{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode request parameters")
		return
	}

	pwHash, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password")
	}

	newUser, err := cfg.DB.CreateUser(request.Email, pwHash)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "user with that password already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "unable to create user")
		return
	}

	newUserResp := response{
		ID:    newUser.ID,
		Email: newUser.Email,
		Red:   newUser.Red,
	}

	respondWithJSON(w, http.StatusCreated, newUserResp)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		respondWithError(w, http.StatusUnauthorized, "missing authorization token")
		return
	}

	userID, err := auth.CheckAuthToken(cfg.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	request := request{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode request parameters")
		return
	}

	pwHash, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password")
		return
	}

	updatedUser, err := cfg.DB.UpdateUserInfo(userID, request.Email, pwHash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update user")
		return
	}

	updatedUserResp := response{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
		Red:   updatedUser.Red,
	}

	respondWithJSON(w, http.StatusOK, updatedUserResp)
}
