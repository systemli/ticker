package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuthenticate(t *testing.T) {
	user, err := NewUser("louis@systemli.org", "password")
	assert.Nil(t, err)

	assert.False(t, user.Authenticate("wrong"))
	assert.True(t, user.Authenticate("password"))
}

func TestUserUpdatePassword(t *testing.T) {
	user, err := NewUser("louis@systemli.org", "password")
	assert.Nil(t, err)

	oldEncPassword := user.EncryptedPassword
	user.UpdatePassword("newPassword")
	assert.NotEqual(t, oldEncPassword, user.EncryptedPassword)
}
