package storage

import (
	. "github.com/systemli/ticker/internal/model"
)

//
func UserExists(data interface{}) bool {
	var user User

	err := DB.One("ID", int(data.(float64)), &user)
	if err != nil {
		return false
	}

	return true
}

//
func UserAuthenticate(email, password string) (*User, bool) {
	var user User

	err := DB.One("Email", email, &user)
	if err != nil {
		return &user, false
	}

	return &user, user.Authenticate(password)
}
