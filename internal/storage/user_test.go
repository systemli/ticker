package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

func TestUserExists(t *testing.T) {
	setup()

	u, err := NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	DB.Save(u)

	assert.True(t, UserExists(u.Email))
	assert.False(t, UserExists("99"))
}

func TestUserAuthenticate(t *testing.T) {
	setup()

	u, err := NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	DB.Save(u)

	user, auth := UserAuthenticate("louis@systemli.org", "password")
	assert.Equal(t, u.ID, user.ID)
	assert.True(t, auth)

	user, auth = UserAuthenticate("louis@systemli.org", "wrong")
	assert.Equal(t, u.ID, user.ID)
	assert.False(t, auth)

	user, auth = UserAuthenticate("admin@systemli.org", "password")
	assert.Equal(t, 0, user.ID)
	assert.False(t, auth)
}
