package user

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPrefetchUserMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	mw := PrefetchUser(s)

	mw(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrefetchUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("userID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("storage error"))
	mw := PrefetchUser(s)

	mw(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrefetchUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("userID", "1")
	s := &storage.MockTickerStorage{}
	user := storage.User{ID: 1}
	s.On("FindUserByID", mock.Anything).Return(user, nil)
	mw := PrefetchUser(s)

	mw(c)

	us, e := c.Get("user")
	assert.True(t, e)
	assert.Equal(t, user, us.(storage.User))
}
