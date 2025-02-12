package response

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type UsersResponseTestSuite struct {
	suite.Suite
}

func (s *UsersResponseTestSuite) TestUsersResponse() {
	users := []storage.User{
		{
			ID:           1,
			CreatedAt:    time.Now(),
			LastLogin:    time.Now(),
			Email:        "user@systemli.org",
			IsSuperAdmin: true,
			Tickers: []storage.Ticker{
				{
					ID:     1,
					Domain: "example.com",
					Title:  "Example",
				},
			},
		},
	}

	usersResponse := UsersResponse(users)
	s.Equal(1, len(usersResponse))
	s.Equal(users[0].ID, usersResponse[0].ID)
	s.Equal(users[0].CreatedAt, usersResponse[0].CreatedAt)
	s.Equal(users[0].LastLogin, usersResponse[0].LastLogin)
	s.Equal(users[0].Email, usersResponse[0].Email)
	s.Equal(users[0].IsSuperAdmin, usersResponse[0].IsSuperAdmin)
	s.Equal(1, len(usersResponse[0].Tickers))
	s.Equal(users[0].Tickers[0].ID, usersResponse[0].Tickers[0].ID)
	s.Equal(users[0].Tickers[0].Domain, usersResponse[0].Tickers[0].Domain)
	s.Equal(users[0].Tickers[0].Title, usersResponse[0].Tickers[0].Title)
}

func TestUsersResponseTestSuite(t *testing.T) {
	suite.Run(t, new(UsersResponseTestSuite))
}
