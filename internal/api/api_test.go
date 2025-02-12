package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/logger"
	"github.com/systemli/ticker/internal/storage"
)

type APITestSuite struct {
	cfg    config.Config
	store  *storage.MockStorage
	logger *logrus.Logger
	suite.Suite
}

func (s *APITestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	logrus.SetOutput(io.Discard)

	s.cfg = config.LoadConfig("")
	s.store = &storage.MockStorage{}

	logger := logger.NewLogrus("debug", "text")
	logger.SetOutput(io.Discard)
	s.logger = logger
}

func (s *APITestSuite) TestHealthz() {
	r := API(s.cfg, s.store, s.logger)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.store.AssertExpectations(s.T())
}

func (s *APITestSuite) TestLogin() {
	s.Run("when password is wrong", func() {
		user, err := storage.NewUser("user@systemli.org", "password")
		s.NoError(err)
		s.store.On("FindUserByEmail", mock.Anything).Return(user, nil)
		r := API(s.cfg, s.store, s.logger)

		body := `{"username":"louis@systemli.org","password":"WRONG"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/admin/login", strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var res response.Response
		err = json.Unmarshal(w.Body.Bytes(), &res)
		s.NoError(err)
		s.Equal(http.StatusUnauthorized, w.Code)
		s.Nil(res.Data)
		s.Equal(res.Error.Code, response.CodeBadCredentials)
		s.Equal(res.Error.Message, response.Unauthorized)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when login is successful", func() {
		user, err := storage.NewUser("user@systemli.org", "password")
		s.NoError(err)
		s.store.On("FindUserByEmail", mock.Anything).Return(user, nil)
		s.store.On("SaveUser", mock.Anything).Return(nil)
		r := API(s.cfg, s.store, s.logger)

		body := `{"username":"louis@systemli.org","password":"password"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/admin/login", strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var res struct {
			Code   int       `json:"code"`
			Expire time.Time `json:"expire"`
			Token  string    `json:"token"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &res)
		s.NoError(err)
		s.Equal(http.StatusOK, w.Code)
		s.Equal(http.StatusOK, res.Code)
		s.NotNil(res.Expire)
		s.NotNil(res.Token)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when save user fails", func() {
		user, err := storage.NewUser("user@systemli.org", "password")
		s.NoError(err)
		s.store.On("FindUserByEmail", mock.Anything).Return(user, nil)
		s.store.On("SaveUser", mock.Anything).Return(errors.New("failed to save user"))

		r := API(s.cfg, s.store, s.logger)

		body := `{"username":"louis@systemli.org","password":"password"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/admin/login", strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
	})
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
