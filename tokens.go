package main

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wtwingate/chirpy/internal/database"
)

func (cfg *apiConfig) createNewToken(user database.User, lifetime int) (string, error) {
	if lifetime == 0 || lifetime > 24*60*60 {
		lifetime = 24 * 60 * 60
	}

	startTime := time.Now().UTC()
	endTime := startTime.Add(time.Duration(lifetime) * time.Second)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(startTime),
		ExpiresAt: jwt.NewNumericDate(endTime),
		Subject:   strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
