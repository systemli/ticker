package user

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type UserTestSuite struct {
	suite.Suite
}

func (s *UserTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *UserTestSuite) TestPrefetchUser() {
	s.Run("when param is missing", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		store := &storage.MockStorage{}
		mw := PrefetchUser(store)

		mw(c)

		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("storage returns error", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("userID", "1")
		store := &storage.MockStorage{}
		store.On("FindUserByID", mock.Anything, mock.Anything).Return(storage.User{}, errors.New("storage error"))
		mw := PrefetchUser(store)

		mw(c)

		s.Equal(http.StatusNotFound, w.Code)
	})

	s.Run("storage returns user", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("userID", "1")
		store := &storage.MockStorage{}
		user := storage.User{ID: 1}
		store.On("FindUserByID", mock.Anything, mock.Anything).Return(user, nil)
		mw := PrefetchUser(store)

		mw(c)

		us, e := c.Get("user")
		s.True(e)
		s.Equal(user, us.(storage.User))
	})
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
