package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		log.Println("missing authorization token")
		w.WriteHeader(401)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		log.Printf("authorization error: %v\n", err)
		w.WriteHeader(401)
		return
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		log.Println(err)
		w.WriteHeader(401)
		return
	}

	type parameters struct {
		Email    string
		Password string
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("could not decode parameters: %v\n", err)
		w.WriteHeader(400)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(userID, params.Email, params.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	updatedUserResp := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
	}

	respondWithJSON(w, 200, updatedUserResp)
}
