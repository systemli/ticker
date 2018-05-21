package storage

import (
	"strconv"

	. "git.codecoop.org/systemli/ticker/internal/model"
)

//
func UserExists(userID string) bool {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return false
	}

	var user User
	err = DB.One("ID", id, &user)
	if err != nil {
		return false
	}

	return true
}

//
func UserAuthenticate(email, password string) (string, bool) {
	var user User

	err := DB.One("Email", email, &user)
	if err != nil {
		return "", false
	}

	return strconv.Itoa(user.ID), user.Authenticate(password)
}
