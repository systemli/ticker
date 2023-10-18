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

func TestGetSettingWithoutParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetSetting(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetSettingInactiveSetting(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	c.AddParam("name", storage.SettingInactiveName)
	s := &storage.MockStorage{}
	s.On("GetInactiveSettings").Return(storage.InactiveSettings{})
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
	c.Set("me", storage.User{IsSuperAdmin: true})
	c.AddParam("name", storage.SettingRefreshInterval)
	s := &storage.MockStorage{}
	s.On("GetRefreshIntervalSettings").Return(storage.RefreshIntervalSettings{})
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetSetting(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutInactiveSettingsMissingBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	body := `broken_json`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", strings.NewReader(body))

	s := &storage.MockStorage{}
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
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	setting := storage.DefaultInactiveSettings()
	body, _ := json.Marshal(setting)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockStorage{}
	s.On("SaveInactiveSettings", mock.Anything).Return(errors.New("storage error"))
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
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	setting := storage.DefaultInactiveSettings()
	body, _ := json.Marshal(setting)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", bytes.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockStorage{}
	s.On("SaveInactiveSettings", mock.Anything).Return(nil)
	s.On("GetInactiveSettings").Return(setting)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutInactiveSettings(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutRefreshIntervalSettingsMissingBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	body := `broken_json`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", strings.NewReader(body))

	s := &storage.MockStorage{}
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
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	body := `{"refreshInterval": 10000}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockStorage{}
	s.On("SaveRefreshIntervalSettings", mock.Anything).Return(errors.New("storage error"))
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
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	setting := storage.DefaultRefreshIntervalSettings()
	body := `{"refreshInterval": 10000}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/settings", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockStorage{}
	s.On("SaveRefreshIntervalSettings", mock.Anything).Return(nil)
	s.On("GetRefreshIntervalSettings").Return(setting)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutRefreshInterval(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
