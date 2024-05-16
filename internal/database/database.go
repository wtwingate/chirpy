package database

import (
	"cmp"
	"encoding/json"
	"os"
	"slices"
	"sync"
)

// Establish a database connection and create a new database
// file if one does not exist.
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  new(sync.RWMutex),
	}
	return &db, db.ensureDB()
}

// Create a new chirp and save it to the database.
func (db *DB) NewChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newChirp := Chirp{
		ID:   len(dbStruct.Chirps) + 1,
		Body: body,
	}
	dbStruct.Chirps[newChirp.ID] = newChirp

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

// Return an array of all chirps in the database.
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirpSlice := []Chirp{}
	for _, v := range dbStruct.Chirps {
		chirpSlice = append(chirpSlice, v)
	}

	slices.SortFunc(chirpSlice, func(a, b Chirp) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return chirpSlice, nil
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
	dbStruct := DBStructure{
		Chirps: chirpMap,
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
