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

	assert.True(t, UserExists("louis@systemli.org"))
	assert.False(t, UserExists("admin@systemli.org"))
}

func TestUserAuthenticate(t *testing.T) {
	setup()

	u, err := NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	DB.Save(u)

	email, auth := UserAuthenticate("louis@systemli.org", "password")
	assert.Equal(t, "louis@systemli.org", email)
	assert.True(t, auth)

	email, auth = UserAuthenticate("louis@systemli.org", "wrong")
	assert.Equal(t, "louis@systemli.org", email)
	assert.False(t, auth)

	email, auth = UserAuthenticate("admin@systemli.org", "password")
	assert.Equal(t, "", email)
	assert.False(t, auth)
}
