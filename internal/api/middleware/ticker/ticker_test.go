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
	c.Set("user", storage.User{})
	s := &storage.MockTickerStorage{}
	mw := PrefetchTicker(s)

	mw(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrefetchTickerNoPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: false, Tickers: []int{2}})
	s := &storage.MockTickerStorage{}
	mw := PrefetchTicker(s)

	mw(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPrefetchTickerStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	ticker := storage.Ticker{ID: 1}
	s.On("FindTickerByID", mock.Anything).Return(ticker, nil)
	mw := PrefetchTicker(s)

	mw(c)

	ti, e := c.Get("ticker")
	assert.True(t, e)
	assert.Equal(t, ticker, ti.(storage.Ticker))

}
