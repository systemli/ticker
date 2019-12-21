package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
)

func TestNewUser(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, "louis@systemli.org", user.Email)
	assert.False(t, user.CreationDate.IsZero())
	assert.False(t, user.IsSuperAdmin)
	assert.NotEmpty(t, user.EncryptedPassword)
}

func TestNewAdminUser(t *testing.T) {
	user, err := model.NewAdminUser("admin@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}

	assert.True(t, user.IsSuperAdmin)
}

func TestUser_Authenticate(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}

	assert.True(t, user.Authenticate("password"))
}

func TestUser_AddTicker(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}
	ticker := model.NewTicker()

	user.AddTicker(*ticker)

	assert.Equal(t, 1, len(user.Tickers))
}

func TestUser_RemoveTicker(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}
	ticker := model.NewTicker()
	user.Tickers = []int{ticker.ID}

	assert.Equal(t, 1, len(user.Tickers))

	user.RemoveTicker(*ticker)

	assert.Equal(t, 0, len(user.Tickers))
}

func TestNewUserResponse(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}

	r := model.NewUserResponse(*user)

	assert.Equal(t, "louis@systemli.org", r.Email)
	assert.False(t, r.IsSuperAdmin)
}

func TestNewUsersResponse(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}

	r := model.NewUsersResponse([]model.User{*user})

	assert.Equal(t, 1, len(r))
}

func TestUser_UpdatePassword(t *testing.T) {
	user, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}
	p1 := user.EncryptedPassword

	user.UpdatePassword("password2")

	assert.NotEqual(t, p1, user.EncryptedPassword)
}
