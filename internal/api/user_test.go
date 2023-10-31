package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type UserTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *UserTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *UserTestSuite) Run(name string, subtest func()) {
	s.T().Run(name, func(t *testing.T) {
		s.w = httptest.NewRecorder()
		s.ctx, _ = gin.CreateTestContext(s.w)
		s.store = &storage.MockStorage{}
		s.cfg = config.LoadConfig("")

		subtest()
	})
}

func (s *UserTestSuite) TestGetUsers() {
	s.Run("when storage returns an error", func() {
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindUsers", mock.Anything).Return([]storage.User{}, errors.New("storage error")).Once()
		h := s.handler()
		h.GetUsers(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns users", func() {
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindUsers", mock.Anything).Return([]storage.User{}, nil).Once()
		h := s.handler()
		h.GetUsers(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *UserTestSuite) TestGetUser() {
	s.Run("when user is missing", func() {
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		h := s.handler()
		h.GetUser(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when insufficient permissions", func() {
		s.ctx.Set("user", storage.User{ID: 1})
		s.ctx.Set("me", storage.User{ID: 2, IsSuperAdmin: false})
		h := s.handler()
		h.GetUser(s.ctx)

		s.Equal(http.StatusForbidden, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("returns user", func() {
		s.ctx.Set("user", storage.User{ID: 1})
		s.ctx.Set("me", storage.User{ID: 1, IsSuperAdmin: false})
		h := s.handler()
		h.GetUser(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *UserTestSuite) TestPostUser() {
	s.Run("when body is missing", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", nil)
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		h := s.handler()
		h.PostUser(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when body is invalid", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader("invalid"))
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		h := s.handler()
		h.PostUser(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when password is too long", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(`{"email":"user@systemli.org","password":"swusp-dud-gust-grong-yuz-swuft-plaft-glact-skast-swem-yen-kom-tut-prisp-gont"}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		h := s.handler()
		h.PostUser(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns an error", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(`{"email":"user@systemli.org","password":"password1234"}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("SaveUser", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PostUser(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when save is successful", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(`{"email":"user@systemli.org","password":"password1234"}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("SaveUser", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PostUser(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *UserTestSuite) TestPutUser() {
	s.Run("when user is missing", func() {
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		h := s.handler()
		h.PutUser(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when body is invalid", func() {
		s.ctx.Set("user", storage.User{ID: 1})
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users", strings.NewReader("invalid"))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutUser(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when tickers are missing", func() {
		s.ctx.Set("user", storage.User{ID: 1})
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users", strings.NewReader(`{"email":"louis@systemli.org","password":"password1234","isSuperAdmin":true,"tickers":[{"id":1}]}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveUser", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutUser(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when save is successful", func() {
		s.ctx.Set("user", storage.User{ID: 1})
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users", strings.NewReader(`{"email":"louis@systemli.org","password":"password1234","isSuperAdmin":true,"tickers":[{"id":1}]}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveUser", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PutUser(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *UserTestSuite) TestDeleteUser() {
	s.Run("when user is missing", func() {
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		h := s.handler()
		h.DeleteUser(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when user is self", func() {
		s.ctx.Set("user", storage.User{ID: 1})
		s.ctx.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
		h := s.handler()
		h.DeleteUser(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns an error", func() {
		s.ctx.Set("user", storage.User{ID: 1})
		s.ctx.Set("me", storage.User{ID: 2, IsSuperAdmin: true})
		s.store.On("DeleteUser", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.DeleteUser(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when delete is successful", func() {
		s.ctx.Set("user", storage.User{ID: 1})
		s.ctx.Set("me", storage.User{ID: 2, IsSuperAdmin: true})
		s.store.On("DeleteUser", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.DeleteUser(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *UserTestSuite) TestPutMe() {
	s.Run("when user is missing", func() {
		h := s.handler()
		h.PutMe(s.ctx)

		s.Equal(http.StatusForbidden, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when form is invalid", func() {
		user, _ := storage.NewUser("user@systemli.org", "password1234")
		s.ctx.Set("me", user)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader("invalid"))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutMe(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when password is wrong", func() {
		user, _ := storage.NewUser("user@systemli.org", "password1234")
		s.ctx.Set("me", user)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader(`{"password":"wrongpassword","newPassword":"password5678"}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutMe(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns an error", func() {
		user, _ := storage.NewUser("user@systemli.org", "password1234")
		s.ctx.Set("me", user)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader(`{"password":"password1234","newPassword":"password5678"}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveUser", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutMe(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when save is successful", func() {
		user, _ := storage.NewUser("user@systemli.org", "password1234")
		s.ctx.Set("me", user)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader(`{"password":"password1234","newPassword":"password5678"}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveUser", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PutMe(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *UserTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
