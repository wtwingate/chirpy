package database

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"

	"golang.org/x/crypto/bcrypt"
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

// Create a new user and save it to the database
func (db *DB) NewUser(email string, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	if _, ok := dbStruct.Users[email]; ok {
		return User{}, errors.New("email address is already registered")
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	newUser := User{
		ID:    len(dbStruct.Users) + 1,
		Email: email,
		Hash:  pwHash,
	}
	dbStruct.Users[newUser.Email] = newUser

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) LoginUser(email string, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	loginUser, ok := dbStruct.Users[email]
	if !ok {
		return User{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword(loginUser.Hash, []byte(password))
	if err != nil {
		return User{}, err
	}

	return loginUser, nil
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

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStruct.Chirps[id]
	if !ok {
		errMsg := fmt.Sprintf("error: could not find chirp ID %v", id)
		return Chirp{}, errors.New(errMsg)
	}
	return chirp, nil
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
	userMap := make(map[string]User)
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
