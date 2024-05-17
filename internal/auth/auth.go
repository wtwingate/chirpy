package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeader = errors.New("no auth header in request")

func HashPassword(password string) (string, error) {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwHash), nil
}

func CheckHashPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateNewToken(userID int, jwtSecret string, lifetime int) (string, error) {
	defaultLifetime := 24 * 60 * 60
	if lifetime == 0 || lifetime > defaultLifetime {
		lifetime = defaultLifetime
	}

	duration := time.Duration(lifetime) * time.Second

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(duration)),
		Subject:   strconv.Itoa(userID),
	})
	return token.SignedString([]byte(jwtSecret))
}
