package storage

import (
	"net/http/httptest"
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

func TestNewUserFilter(t *testing.T) {
	filter := NewUserFilter(nil)
	assert.Nil(t, filter.Email)
	assert.Nil(t, filter.IsSuperAdmin)

	req := httptest.NewRequest("GET", "/", nil)
	filter = NewUserFilter(req)
	assert.Nil(t, filter.Email)
	assert.Nil(t, filter.IsSuperAdmin)

	req = httptest.NewRequest("GET", "/?email=user@example.org&is_super_admin=true", nil)
	filter = NewUserFilter(req)
	assert.Equal(t, "user@example.org", *filter.Email)
	assert.True(t, *filter.IsSuperAdmin)

	req = httptest.NewRequest("GET", "/?order_by=created_at&sort=asc", nil)
	filter = NewUserFilter(req)
	assert.Nil(t, filter.Email)
	assert.Nil(t, filter.IsSuperAdmin)
	assert.Equal(t, "created_at", filter.OrderBy)
	assert.Equal(t, "asc", filter.Sort)
}
