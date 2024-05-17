package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wtwingate/chirpy/internal/auth"
	"github.com/wtwingate/chirpy/internal/database"
)

type request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
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
	}

	respondWithJSON(w, http.StatusCreated, newUserResp)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		respondWithError(w, http.StatusUnauthorized, "missing authorization token")
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid authorization token")
		return
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid authorization token")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid user ID")
		return
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

	updatedUser, err := cfg.DB.UpdateUser(userID, request.Email, pwHash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update user")
		return
	}

	updatedUserResp := response{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
	}

	respondWithJSON(w, http.StatusOK, updatedUserResp)
}
