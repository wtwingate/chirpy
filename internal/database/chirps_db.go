package database

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

// Create a new chirp and save it to the database.
func (db *DB) CreateChirp(userID int, body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newChirp := Chirp{
		ID:       len(dbStruct.Chirps) + 1,
		Body:     body,
		AuthorID: userID,
	}
	dbStruct.Chirps[newChirp.ID] = newChirp

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

func (db *DB) DeleteChirp(chirpID, userID int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbStruct.Chirps[chirpID]
	if !ok {
		errMsg := fmt.Sprintf("could not find chirp %v", chirpID)
		return errors.New(errMsg)
	}

	if userID != chirp.AuthorID {
		return errors.New("unauthorized")
	}

	delete(dbStruct.Chirps, chirpID)

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
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
		errMsg := fmt.Sprintf("could not find chirp %v", id)
		return Chirp{}, errors.New(errMsg)
	}
	return chirp, nil
}
