package database

import (
	"errors"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

// Create a new user and save it to the database
func (db *DB) CreateUser(email string, hash string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newUser := User{
		ID:    len(dbStruct.Users) + 1,
		Email: email,
		Hash:  hash,
	}
	dbStruct.Users[newUser.ID] = newUser

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) UpdateUser(id int, email string, hash string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, err := db.GetUser(id)
	if err != nil {
		return User{}, err
	}

	user.Email = email
	user.Hash = hash
	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStruct.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrNotExist
}
