package ticker

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPrefetchTickerParamMissing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{})
	s := &storage.MockTickerStorage{}
	mw := PrefetchTicker(s)

	mw(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrefetchTickerNoPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("me", storage.User{IsSuperAdmin: false, Tickers: []int{2}})
	s := &storage.MockTickerStorage{}
	mw := PrefetchTicker(s)

	mw(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPrefetchTickerStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("me", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, errors.New("storage error"))
	mw := PrefetchTicker(s)

	mw(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrefetchTicker(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("me", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	ticker := storage.Ticker{ID: 1}
	s.On("FindTickerByID", mock.Anything).Return(ticker, nil)
	mw := PrefetchTicker(s)

	mw(c)

	ti, e := c.Get("ticker")
	assert.True(t, e)
	assert.Equal(t, ticker, ti.(storage.Ticker))

}

func TestPrefetchTickerFromRequestMissingOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
	s := &storage.MockTickerStorage{}
	mw := PrefetchTickerFromRequest(s)

	mw(c)

	assert.Equal(t, http.StatusOK, w.Code)
	ticker, exists := c.Get("ticker")
	assert.Equal(t, nil, ticker)
	assert.False(t, exists)
}

func TestPrefetchTickerFromRequestTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
	c.Request.Header.Set("Origin", "https://demoticker.org")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByDomain", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
	mw := PrefetchTickerFromRequest(s)

	mw(c)

	assert.Equal(t, http.StatusOK, w.Code)
	ticker, exists := c.Get("ticker")
	assert.Equal(t, nil, ticker)
	assert.False(t, exists)
}

func TestPrefetchTickerFromRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
	c.Request.Header.Set("Origin", "https://demoticker.org")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByDomain", mock.Anything).Return(storage.Ticker{}, nil)
	mw := PrefetchTickerFromRequest(s)

	mw(c)

	assert.Equal(t, http.StatusOK, w.Code)
	ticker, exists := c.Get("ticker")
	assert.NotNil(t, ticker)
	assert.True(t, exists)
}
