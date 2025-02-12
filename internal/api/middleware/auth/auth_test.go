package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

type AuthTestSuite struct {
	suite.Suite
}

func (s *AuthTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *AuthTestSuite) TestAuthenticator() {
	s.Run("when form is empty", func() {
		mockStorage := &storage.MockStorage{}
		authenticator := Authenticator(mockStorage)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{}`))

		_, err := authenticator(c)
		s.Error(err)
		s.Equal("missing Username or Password", err.Error())
	})

	s.Run("when user is not found", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUserByEmail", mock.Anything).Return(storage.User{}, errors.New("not found"))

		authenticator := Authenticator(mockStorage)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "user@systemli.org", "password": "password"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		_, err := authenticator(c)
		s.Error(err)
		s.Equal("not found", err.Error())
	})

	s.Run("when user is found", func() {
		user, err := storage.NewUser("user@systemli.org", "password")
		s.NoError(err)

		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUserByEmail", mock.Anything).Return(user, nil)
		mockStorage.On("SaveUser", mock.Anything).Return(nil)
		authenticator := Authenticator(mockStorage)

		s.Run("with correct password", func() {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "user@systemli.org", "password": "password"}`))
			c.Request.Header.Set("Content-Type", "application/json")
			user, err := authenticator(c)

			s.NoError(err)
			s.Equal("user@systemli.org", user.(storage.User).Email)
		})

		s.Run("with incorrect password", func() {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "user@systemli.org", "password": "password1"}`))
			c.Request.Header.Set("Content-Type", "application/json")

			_, err := authenticator(c)

			s.Error(err)
			s.Equal("authentication failed", err.Error())
		})
	})
}

func (s *AuthTestSuite) TestAuthorizator() {
	s.Run("when user is not found", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("not found"))
		authorizator := Authorizator(mockStorage)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		found := authorizator(float64(1), c)
		s.False(found)
	})

	s.Run("when user is found", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUserByID", mock.Anything).Return(storage.User{ID: 1}, nil)
		authorizator := Authorizator(mockStorage)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		found := authorizator(float64(1), c)
		s.True(found)
	})
}

func (s *AuthTestSuite) TestUnauthorized() {
	s.Run("returns a 403 with json payload", func() {
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)
		c.Request = httptest.NewRequest(http.MethodGet, "/login", nil)

		Unauthorized(c, 403, "unauthorized")

		err := json.Unmarshal(rr.Body.Bytes(), &response.Response{})
		s.NoError(err)
		s.Equal(403, rr.Code)
	})
}

func (s *AuthTestSuite) TestFillClaims() {
	s.Run("when user is empty", func() {
		claims := FillClaim("empty")
		s.Equal(jwt.MapClaims{}, claims)
	})

	s.Run("when user is valid", func() {
		user := storage.User{ID: 1, Email: "user@systemli.org", IsSuperAdmin: true}
		claims := FillClaim(user)

		s.Equal(jwt.MapClaims{"id": 1, "email": "user@systemli.org", "roles": []string{"user", "admin"}}, claims)
	})
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
