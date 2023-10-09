package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Password        = "password"
	TooLongPassword = "swusp-dud-gust-grong-yuz-swuft-plaft-glact-skast-swem-yen-kom-tut-prisp-gont"
)

func TestUserAuthenticate(t *testing.T) {
	user, err := NewUser("louis@systemli.org", Password)
	assert.Nil(t, err)

	assert.False(t, user.Authenticate("wrong"))
	assert.True(t, user.Authenticate(Password))
}

func TestUserUpdatePassword(t *testing.T) {
	user, err := NewUser("louis@systemli.org", Password)
	assert.Nil(t, err)

	oldEncPassword := user.EncryptedPassword
	user.UpdatePassword("newPassword")
	assert.NotEqual(t, oldEncPassword, user.EncryptedPassword)

	user.UpdatePassword(TooLongPassword)
	assert.NotEqual(t, oldEncPassword, user.EncryptedPassword)
}

func TestNewUser(t *testing.T) {
	_, err := NewUser("user@systemli.org", Password)
	assert.Nil(t, err)

	_, err = NewUser("user@systemli.org", TooLongPassword)
	assert.NotNil(t, err)
}
