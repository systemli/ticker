package response

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/storage"
)

func TestUsersResponse(t *testing.T) {
	users := []storage.User{
		{
			ID:           1,
			CreatedAt:    time.Now(),
			Email:        "user@systemli.org",
			IsSuperAdmin: true,
		},
	}

	usersResponse := UsersResponse(users)
	assert.Equal(t, 1, len(usersResponse))
	assert.Equal(t, users[0].ID, usersResponse[0].ID)
	assert.Equal(t, users[0].CreatedAt, usersResponse[0].CreatedAt)
	assert.Equal(t, users[0].Email, usersResponse[0].Email)
	assert.Equal(t, users[0].IsSuperAdmin, usersResponse[0].IsSuperAdmin)
}
