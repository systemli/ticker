package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

func TestFindUserByID(t *testing.T) {
	setup()

	u, err := NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	DB.Save(u)

	user, err := FindUserByID(u.ID)

	assert.Equal(t, u.ID, user.ID)
	assert.Nil(t, err)

	user, err = FindUserByID(2)

	assert.Equal(t, 0, user.ID)
	assert.NotNil(t, err)
}

func TestUserAuthenticate(t *testing.T) {
	setup()

	u, err := NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	DB.Save(u)

	user, err := UserAuthenticate("louis@systemli.org", "password")
	assert.Equal(t, u.ID, user.ID)
	assert.Nil(t, err)

	user, err = UserAuthenticate("louis@systemli.org", "wrong")
	assert.Equal(t, u.ID, user.ID)
	assert.NotNil(t, err)

	user, err = UserAuthenticate("admin@systemli.org", "password")
	assert.Equal(t, 0, user.ID)
	assert.NotNil(t, err)
}
