package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetTimelineMissingDomain(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTimeline(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), response.TickerNotFound)
}

func TestGetTimelineTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
	c.Request.Header.Add("Origin", "https://demoticker.org")

	s := &storage.MockTickerStorage{}
	s.On("FindTickerByDomain", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTimeline(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), response.TickerNotFound)
}

func TestGetTimelineMessageFetchError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	s.On("FindMessagesByTickerAndPagination", mock.Anything, mock.Anything).Return([]storage.Message{}, errors.New("storage error"))

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTimeline(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), response.MessageFetchError)
}

func TestGetTimeline(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	s.On("FindMessagesByTickerAndPagination", mock.Anything, mock.Anything).Return([]storage.Message{}, nil)

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTimeline(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
