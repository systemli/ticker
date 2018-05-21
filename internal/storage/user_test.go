package storage_test

import (
	"strconv"
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

	assert.True(t, UserExists(strconv.Itoa(u.ID)))
	assert.False(t, UserExists("99"))
}

func TestUserAuthenticate(t *testing.T) {
	setup()

	u, err := NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	DB.Save(u)

	id, auth := UserAuthenticate("louis@systemli.org", "password")
	assert.Equal(t, strconv.Itoa(u.ID), id)
	assert.True(t, auth)

	id, auth = UserAuthenticate("louis@systemli.org", "wrong")
	assert.Equal(t, strconv.Itoa(u.ID), id)
	assert.False(t, auth)

	id, auth = UserAuthenticate("admin@systemli.org", "password")
	assert.Equal(t, "", id)
	assert.False(t, auth)
}
