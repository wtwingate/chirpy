package database

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func (db *DB) CreateNewRefreshToken() (string, error) {
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
	dbStruct.Refresh[token] = time.Now().UTC().Add(60 * 24 * time.Hour)

	err = db.writeDB(dbStruct)
	if err != nil {
		return "", err
	}

	return token, nil
}
