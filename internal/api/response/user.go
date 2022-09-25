package response

import (
	"time"

	"github.com/systemli/ticker/internal/storage"
)

type User struct {
	ID           int       `json:"id"`
	CreationDate time.Time `json:"creation_date"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	Tickers      []int     `json:"tickers"`
}

func UserResponse(user storage.User) User {
	return User{
		ID:           user.ID,
		CreationDate: user.CreationDate,
		Email:        user.Email,
		Role:         user.Role,
		IsSuperAdmin: user.IsSuperAdmin,
		Tickers:      user.Tickers,
	}
}

func UsersResponse(users []storage.User) []User {
	u := make([]User, 0)
	for _, user := range users {
		u = append(u, UserResponse(user))
	}

	return u
}
