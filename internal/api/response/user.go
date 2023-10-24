package response

import (
	"time"

	"github.com/systemli/ticker/internal/storage"
)

type User struct {
	ID           int          `json:"id"`
	CreatedAt    time.Time    `json:"createdAt"`
	Email        string       `json:"email"`
	Role         string       `json:"role"`
	Tickers      []UserTicker `json:"tickers"`
	IsSuperAdmin bool         `json:"isSuperAdmin"`
}

type UserTicker struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Title  string `json:"title"`
}

func UserResponse(user storage.User) User {
	return User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		Email:        user.Email,
		IsSuperAdmin: user.IsSuperAdmin,
		Tickers:      UserTickersResponse(user.Tickers),
	}
}

func UsersResponse(users []storage.User) []User {
	u := make([]User, 0)
	for _, user := range users {
		u = append(u, UserResponse(user))
	}

	return u
}

func UserTickersResponse(tickers []storage.Ticker) []UserTicker {
	t := make([]UserTicker, 0)
	for _, ticker := range tickers {
		t = append(t, UserTickerResponse(ticker))
	}

	return t
}

func UserTickerResponse(ticker storage.Ticker) UserTicker {
	return UserTicker{
		ID:     ticker.ID,
		Domain: ticker.Domain,
		Title:  ticker.Title,
	}
}
