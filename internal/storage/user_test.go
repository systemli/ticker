package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

func TestFindUserByID(t *testing.T) {
	setup()

	u := initUserTestData(t)

	user, err := FindUserByID(u.ID)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, u.ID, user.ID)
	assert.Nil(t, err)

	_, err = FindUserByID(2)
	assert.NotNil(t, err)
}

func TestFindUsers(t *testing.T) {
	setup()

	users, err := FindUsers()
	if err == nil {
		t.Fail()
		return
	}

	u := initUserTestData(t)

	users, err = FindUsers()
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 1, len(users))
	assert.Equal(t, u.ID, users[0].ID)
}

func TestFindUsersByTicker(t *testing.T) {
	setup()

	ticker := NewTicker()
	_ = DB.Save(ticker)

	users, err := FindUsersByTicker(*ticker)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 0, len(users))

	u := initUserTestData(t)
	u.Tickers = []int{ticker.ID}
	_ = DB.Save(u)

	users, err = FindUsersByTicker(*ticker)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 1, len(users))
}

func TestUserAuthenticate(t *testing.T) {
	setup()

	u := initUserTestData(t)

	user, err := UserAuthenticate("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, u.ID, user.ID)

	user, err = UserAuthenticate("louis@systemli.org", "wrong")
	assert.NotNil(t, err)

	user, err = UserAuthenticate("admin@systemli.org", "password")
	assert.NotNil(t, err)
}

func TestAddUsersToTicker(t *testing.T) {
	setup()

	u := initUserTestData(t)
	ticker := NewTicker()
	_ = DB.Save(ticker)

	err := AddUsersToTicker(*ticker, []int{u.ID})
	if err != nil {
		t.Fail()
	}

	var user User
	err = DB.One("ID", 1, &user)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 1, len(user.Tickers))

	err = AddUsersToTicker(*ticker, []int{2})
	if err == nil {
		t.Fail()
	}

	admin, err := NewAdminUser("admin@systemli.org", "password")
	if err != nil {
		t.Fail()
		return
	}
	_ = DB.Save(admin)

	err = AddUsersToTicker(*ticker, []int{admin.ID})
	if err != nil {
		t.Fail()
	}
}

func TestRemoveTickerFromUser(t *testing.T) {
	setup()

	user := initUserTestData(t)
	ticker := NewTicker()
	_ = DB.Save(ticker)
	user.Tickers = []int{ticker.ID}
	_ = DB.Save(user)

	assert.Equal(t, 1, len(user.Tickers))

	err := RemoveTickerFromUser(*ticker, *user)
	if err != nil {
		t.Fail()
		return
	}

	err = DB.One("ID", 1, user)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, len(user.Tickers))
}

func initUserTestData(t *testing.T) *User {
	u, err := NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	_ = DB.Save(u)

	return u
}
