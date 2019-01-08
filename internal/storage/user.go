package storage

import (
	"errors"
	. "github.com/systemli/ticker/internal/model"
)

//FindUserByID returns user if one exists with the given id.
func FindUserByID(id int) (*User, error) {
	var user User

	err := DB.One("ID", id, &user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

//UserAuthenticate returns User when authentication was successful.
func UserAuthenticate(email, password string) (*User, error) {
	var user User

	err := DB.One("Email", email, &user)
	if err != nil {
		return &user, err
	}

	if user.Authenticate(password) {
		return &user, nil
	}

	return &user, errors.New("authentication failed")
}
