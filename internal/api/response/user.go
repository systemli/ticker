package response

import (
	"time"

	"github.com/systemli/ticker/internal/storage"
)

type User struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	IsSuperAdmin bool      `json:"isSuperAdmin"`
}

func UserResponse(user storage.User) User {
	return User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		Email:        user.Email,
		IsSuperAdmin: user.IsSuperAdmin,
	}
}

func UsersResponse(users []storage.User) []User {
	u := make([]User, 0)
	for _, user := range users {
		u = append(u, UserResponse(user))
	}

	return u
}
