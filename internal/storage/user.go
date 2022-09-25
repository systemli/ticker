package storage

import (
	"time"

	"github.com/systemli/ticker/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int       `storm:"id,increment"`
	CreationDate      time.Time `storm:"index"`
	Email             string    `storm:"unique"`
	Role              string
	EncryptedPassword string
	IsSuperAdmin      bool
	Tickers           []int
}

func NewUser(email, password string) (User, error) {
	user := User{
		CreationDate:      time.Now(),
		IsSuperAdmin:      false,
		Email:             email,
		Tickers:           []int{},
		EncryptedPassword: "",
		Role:              "",
	}

	pw, err := hashPassword(password)
	if err != nil {
		return user, err
	}

	user.EncryptedPassword = pw

	return user, nil
}

func NewAdminUser(email, password string) (User, error) {
	user, err := NewUser(email, password)
	if err != nil {
		return user, err
	}

	user.IsSuperAdmin = true

	return user, err
}

func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
}

func (u *User) AddTicker(ticker Ticker) {
	u.Tickers = util.Append(u.Tickers, ticker.ID)
}

func (u *User) RemoveTicker(ticker Ticker) {
	u.Tickers = util.Remove(u.Tickers, ticker.ID)
}

func (u *User) UpdatePassword(password string) {
	pw, err := hashPassword(password)
	if err != nil {
		return
	}

	u.EncryptedPassword = pw
}

func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}
