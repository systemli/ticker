package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

//
type User struct {
	ID                int       `storm:"id,increment"`
	CreationDate      time.Time `storm:"index"`
	Email             string    `storm:"unique"`
	Role              string
	EncryptedPassword string
	IsSuperAdmin      bool
	Tickers           []int
}

//
type UserResponse struct {
	ID           int       `json:"id"`
	CreationDate time.Time `json:"creation_date"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	Tickers      []int     `json:"tickers"`
}

//NewUser returns a new User.
func NewUser(email, password string) (*User, error) {
	pw, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		CreationDate:      time.Now(),
		Email:             email,
		EncryptedPassword: pw,
	}
	return user, nil
}

//NewAdminUser returns a Admin User.
func NewAdminUser(email, password string) (*User, error) {
	user, err := NewUser(email, password)
	user.IsSuperAdmin = true

	return user, err
}

//
func NewUserResponse(user User) *UserResponse {
	return &UserResponse{
		ID:           user.ID,
		CreationDate: user.CreationDate,
		Email:        user.Email,
		Role:         user.Role,
		IsSuperAdmin: user.IsSuperAdmin,
		Tickers:      user.Tickers,
	}
}

func NewUsersResponse(users []User) []*UserResponse {
	var u []*UserResponse

	for _, user := range users {
		u = append(u, NewUserResponse(user))
	}

	return u
}

//
func (u *User) UpdatePassword(password string) {
	pw, err := hashPassword(password)
	if err != nil {
		return
	}

	u.EncryptedPassword = pw
}

// Authenticate a user from a password
func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
}

// hashPassword generates a hashed password from a plaintext string
func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}
