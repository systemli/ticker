package storage

import . "git.codecoop.org/systemli/ticker/internal/model"

//
func UserExists(email string) bool {
	var user User

	err := DB.One("Email", email, &user)
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

	return user.Email, user.Authenticate(password)
}
