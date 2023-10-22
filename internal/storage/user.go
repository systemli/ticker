package storage

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Email             string `storm:"unique"`
	EncryptedPassword string
	IsSuperAdmin      bool
	Tickers           []Ticker `gorm:"many2many:ticker_users;"`
}

func NewUser(email, password string) (User, error) {
	user := User{
		IsSuperAdmin:      false,
		Email:             email,
		EncryptedPassword: "",
	}

	pw, err := hashPassword(password)
	if err != nil {
		return user, err
	}

	user.EncryptedPassword = pw

	return user, nil
}

func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
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
