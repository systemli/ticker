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

	s.Run("get telegram settings", func() {
		s.ctx.AddParam("name", storage.SettingTelegramName)
		expectedSettings := storage.TelegramSettings{Token: "123456789:ABCdefGHIjklMNOpqrsTUVwxyz", BotUsername: "test_bot"}
		s.store.On("GetTelegramSettings").Return(expectedSettings).Once()
		h := s.handler()
		h.GetSetting(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), "telegram_settings")
		// Token should be masked in response
		s.NotContains(s.w.Body.String(), "123456789:ABCdefGHIjklMNOpqrsTUVwxyz")
		s.Contains(s.w.Body.String(), "****wxyz")
		s.Contains(s.w.Body.String(), "test_bot")
		s.store.AssertExpectations(s.T())
	})

	s.Run("get signal group settings", func() {
		s.ctx.AddParam("name", storage.SettingSignalGroupName)
		expectedSettings := storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
			Avatar:  "/path/to/avatar.png",
		}
		s.store.On("GetSignalGroupSettings").Return(expectedSettings).Once()
		h := s.handler()
		h.GetSetting(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), "signal_group_settings")
		s.Contains(s.w.Body.String(), "https://signal-cli.example.org/api/v1/rpc")
		s.Contains(s.w.Body.String(), "0123456789")
		s.Contains(s.w.Body.String(), "/path/to/avatar.png")
		s.store.AssertExpectations(s.T())
	})

	s.Run("get unknown setting", func() {
		s.ctx.AddParam("name", "unknown_setting")
		h := s.handler()
		h.GetSetting(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
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

func (s *SettingsTestSuite) TestPutTelegramSettings() {
	s.Run("when body is invalid", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/settings/telegram_settings", strings.NewReader(`{"token":123}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutTelegramSettings(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		setting := storage.TelegramSettings{Token: ""}
		body, _ := json.Marshal(setting)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/settings/telegram_settings", bytes.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTelegramSettings", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutTelegramSettings(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when saving empty token", func() {
		setting := storage.TelegramSettings{Token: ""}
		body, _ := json.Marshal(setting)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/settings/telegram_settings", bytes.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTelegramSettings", mock.Anything).Return(nil).Once()
		s.store.On("GetTelegramSettings").Return(setting)
		h := s.handler()
		h.PutTelegramSettings(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), "telegram_settings")
		s.NotContains(s.w.Body.String(), "ABCdefGHI")
		s.store.AssertExpectations(s.T())
	})
}

func (s *SettingsTestSuite) TestPutSignalGroupSettings() {
	s.Run("when body is invalid", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/settings/signal_group_settings", strings.NewReader(`{"apiUrl":123}`))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutSignalGroupSettings(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		setting := storage.SignalGroupSettings{ApiUrl: "", Account: ""}
		body, _ := json.Marshal(setting)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/settings/signal_group_settings", bytes.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveSignalGroupSettings", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutSignalGroupSettings(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when saving empty settings", func() {
		setting := storage.SignalGroupSettings{ApiUrl: "", Account: "", Avatar: ""}
		body, _ := json.Marshal(setting)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/settings/signal_group_settings", bytes.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveSignalGroupSettings", mock.Anything).Return(nil).Once()
		s.store.On("GetSignalGroupSettings").Return(setting)
		h := s.handler()
		h.PutSignalGroupSettings(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), "signal_group_settings")
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
