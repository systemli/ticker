package model_test

import (
	"testing"
	"git.codecoop.org/systemli/ticker/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, "louis@systemli.org", user.Email)
	assert.True(t, !user.CreationDate.IsZero())
	assert.NotEmpty(t, user.EncryptedPassword)
}

func TestUser_Authenticate(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	assert.True(t, user.Authenticate("password"))
}