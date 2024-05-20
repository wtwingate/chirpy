package auth

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeader = errors.New("no auth header in request")
var ErrInvalidAuthToken = errors.New("invalid authentication token")

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

func CreateNewAuthToken(userID int, jwtSecret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		Subject:   strconv.Itoa(userID),
	})
	return token.SignedString([]byte(jwtSecret))
}

func CheckAuthToken(secret, token string) (int, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	claims := jwt.RegisteredClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, ErrInvalidAuthToken
	}

	subject, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return 0, ErrInvalidAuthToken
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		return 0, ErrInvalidAuthToken
	}

	return userID, nil
}
