package database

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Refresh struct {
	UserID    int
	Token     string
	ExpiresAt time.Time
}

func (db *DB) CreateNewRefreshToken(userID int) (string, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return "", err
	}

	randBytes := make([]byte, 32)
	_, err = rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(randBytes)
	refresh := Refresh{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	}
	dbStruct.Refresh[token] = refresh

	err = db.writeDB(dbStruct)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (db *DB) CheckRefreshToken(token string) (int, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return 0, err
	}

	refresh, ok := dbStruct.Refresh[token]
	if !ok {
		return 0, ErrNotExist
	}

	return refresh.UserID, nil
}

func (db *DB) RevokeRefreshToken(token string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStruct.Refresh[token]
	if !ok {
		return ErrNotExist
	}

	delete(dbStruct.Refresh, token)

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}
