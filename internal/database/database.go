package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var ErrNotExist = errors.New("resource does not exist")
var ErrAlreadyExists = errors.New("resource already exists")

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps  map[int]Chirp
	Users   map[int]User
	Refresh map[string]time.Time
}

// Establish a database connection and create a new database
// file if one does not exist.
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  new(sync.RWMutex),
	}
	return &db, db.ensureDB()
}

// Create a new database if one does not exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// Read the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	chirpMap := make(map[int]Chirp)
	userMap := make(map[int]User)
	dbStruct := DBStructure{
		Chirps: chirpMap,
		Users:  userMap,
	}

	if len(data) == 0 {
		return dbStruct, nil
	}

	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return DBStructure{}, err
	}
	return dbStruct, nil
}

// Write database structure to disk
func (db *DB) writeDB(dbStruct DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStruct)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}
