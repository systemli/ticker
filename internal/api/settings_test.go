package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type SettingsTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *SettingsTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *SettingsTestSuite) Run(name string, subtest func()) {
	s.T().Run(name, func(t *testing.T) {
		s.w = httptest.NewRecorder()
		s.ctx, _ = gin.CreateTestContext(s.w)
		s.ctx.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
		s.store = &storage.MockStorage{}
		s.cfg = config.LoadConfig("")

		subtest()
	})
}

func (s *SettingsTestSuite) TestGetSetting() {
	s.Run("when url param is missing", func() {
		h := s.handler()
		h.GetSetting(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("get inactive settings", func() {
		s.ctx.AddParam("name", storage.SettingInactiveName)
		s.store.On("GetInactiveSettings").Return(storage.InactiveSettings{}).Once()
		h := s.handler()
		h.GetSetting(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("get refresh interval settings", func() {
		s.ctx.AddParam("name", storage.SettingRefreshInterval)
		s.store.On("GetRefreshIntervalSettings").Return(storage.RefreshIntervalSettings{}).Once()
		h := s.handler()
		h.GetSetting(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when user is no admin", func() {
		s.ctx.Set("me", storage.User{ID: 1, IsSuperAdmin: false})
		s.ctx.AddParam("name", storage.SettingRefreshInterval)
		h := s.handler()
		h.GetSetting(s.ctx)

		s.Equal(http.StatusForbidden, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *SettingsTestSuite) TestPutInactiveSetting() {
	s.Run("when body is invalid", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", strings.NewReader(`broken_json`))
		h := s.handler()
		h.PutInactiveSettings(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		setting := storage.DefaultInactiveSettings()
		body, _ := json.Marshal(setting)
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveInactiveSettings", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutInactiveSettings(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns settings", func() {
		setting := storage.DefaultInactiveSettings()
		body, _ := json.Marshal(setting)
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveInactiveSettings", mock.Anything).Return(nil).Once()
		s.store.On("GetInactiveSettings").Return(setting)
		h := s.handler()
		h.PutInactiveSettings(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *SettingsTestSuite) TestPutRefreshIntervalSetting() {
	s.Run("when body is invalid", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", strings.NewReader(`broken_json`))
		h := s.handler()
		h.PutRefreshInterval(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		setting := storage.DefaultRefreshIntervalSettings()
		body, _ := json.Marshal(setting)
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveRefreshIntervalSettings", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutRefreshInterval(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns settings", func() {
		setting := storage.DefaultRefreshIntervalSettings()
		body, _ := json.Marshal(setting)
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveRefreshIntervalSettings", mock.Anything).Return(nil).Once()
		s.store.On("GetRefreshIntervalSettings").Return(setting)
		h := s.handler()
		h.PutRefreshInterval(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *SettingsTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestSettingsTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsTestSuite))
}
