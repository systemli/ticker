package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func TestGetInit(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/init?origin=demoticker.org", nil)

	ticker := storage.NewTicker()
	ticker.Active = true
	s := &storage.MockStorage{}
	s.On("GetRefreshIntervalSettings").Return(storage.DefaultRefreshIntervalSettings())
	s.On("FindTickerByDomain", mock.AnythingOfType("string")).Return(ticker, nil)

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetInit(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetInitInvalidDomain(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/init", nil)

	s := &storage.MockStorage{}
	s.On("GetRefreshIntervalSettings").Return(storage.DefaultRefreshIntervalSettings())

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetInit(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetInitInactiveTicker(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/init?origin=demoticker.org", nil)

	ticker := storage.NewTicker()
	s := &storage.MockStorage{}
	s.On("GetRefreshIntervalSettings").Return(storage.DefaultRefreshIntervalSettings())
	s.On("GetInactiveSettings").Return(storage.DefaultInactiveSettings())
	s.On("FindTickerByDomain", mock.AnythingOfType("string")).Return(ticker, nil)

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetInit(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
