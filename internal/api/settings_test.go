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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetSettingForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetSetting(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetSettingWithoutParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetSetting(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSettingInactiveSetting(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("name", storage.SettingInactiveName)
	s := &storage.MockTickerStorage{}
	s.On("GetInactiveSetting").Return(storage.Setting{})
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetSetting(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSettingRefreshIntervalSetting(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("name", storage.SettingRefreshInterval)
	s := &storage.MockTickerStorage{}
	s.On("GetRefreshIntervalSetting").Return(storage.Setting{})
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetSetting(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutInactiveSettingsForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutInactiveSettings(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPutInactiveSettingsMissingBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	body := `broken_json`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", strings.NewReader(body))

	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutInactiveSettings(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutInactiveSettingsStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	setting := storage.DefaultInactiveSetting()
	body, _ := json.Marshal(setting.Value)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockTickerStorage{}
	s.On("SaveInactiveSetting", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutInactiveSettings(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutInactiveSettings(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	setting := storage.DefaultInactiveSetting()
	body, _ := json.Marshal(setting.Value)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockTickerStorage{}
	s.On("SaveInactiveSetting", mock.Anything).Return(nil)
	s.On("GetInactiveSetting").Return(setting)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutInactiveSettings(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutRefreshIntervalSettingsForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutRefreshInterval(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPutRefreshIntervalSettingsMissingBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	body := `broken_json`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", strings.NewReader(body))

	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutRefreshInterval(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutRefreshIntervalSettingsStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	setting := storage.DefaultRefreshIntervalSetting()
	body, _ := json.Marshal(setting.Value)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockTickerStorage{}
	s.On("SaveRefreshInterval", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutRefreshInterval(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutRefreshIntervalSettings(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	setting := storage.DefaultRefreshIntervalSetting()
	body, _ := json.Marshal(setting.Value)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockTickerStorage{}
	s.On("SaveRefreshInterval", mock.Anything).Return(nil)
	s.On("GetRefreshIntervalSetting").Return(setting)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutRefreshInterval(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
