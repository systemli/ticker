package me

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

type MeTestSuite struct {
	suite.Suite
}

func (s *MeTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *MeTestSuite) TestMeMiddleware() {
	s.Run("when id is not present", func() {
		mockStorage := &storage.MockStorage{}
		mw := MeMiddleware(mockStorage)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mw(c)

		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("when id is present", func() {
		s.Run("when user is not found", func() {
			mockStorage := &storage.MockStorage{}
			mockStorage.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("not found"))
			mw := MeMiddleware(mockStorage)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("id", float64(1))

			mw(c)

			s.Equal(http.StatusBadRequest, w.Code)
		})

		s.Run("when user is found", func() {
			mockStorage := &storage.MockStorage{}
			mockStorage.On("FindUserByID", mock.Anything).Return(storage.User{ID: 1}, nil)
			mw := MeMiddleware(mockStorage)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("id", float64(1))

			mw(c)

			user, exists := c.Get("me")
			s.True(exists)
			s.IsType(storage.User{}, user)
		})
	})
}

func TestMeTestSuite(t *testing.T) {
	suite.Run(t, new(MeTestSuite))
}
