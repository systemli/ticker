package ticker

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type TickerTestSuite struct {
	suite.Suite
}

func (s *TickerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *TickerTestSuite) TestPrefetchTicker() {
	s.Run("when param is missing", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		store := &storage.MockStorage{}
		mw := PrefetchTicker(store)

		mw(c)

		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("storage returns error", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("tickerID", "1")
		store := &storage.MockStorage{}
		store.On("FindTickerByUserAndID", mock.Anything, mock.Anything, mock.Anything).Return(storage.Ticker{}, errors.New("storage error"))
		mw := PrefetchTicker(store)

		mw(c)

		s.Equal(http.StatusNotFound, w.Code)
	})

	s.Run("storage returns ticker", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("tickerID", "1")
		store := &storage.MockStorage{}
		ticker := storage.Ticker{ID: 1}
		store.On("FindTickerByUserAndID", mock.Anything, mock.Anything, mock.Anything).Return(ticker, nil)
		mw := PrefetchTicker(store)

		mw(c)

		ti, e := c.Get("ticker")
		s.True(e)
		s.Equal(ticker, ti.(storage.Ticker))
	})
}

func (s *TickerTestSuite) TestPrefetchTickerFromRequest() {
	s.Run("when origin is missing", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
		store := &storage.MockStorage{}
		mw := PrefetchTickerFromRequest(store)

		mw(c)

		s.Equal(http.StatusOK, w.Code)
		ticker, exists := c.Get("ticker")
		s.Nil(ticker)
		s.False(exists)
	})

	s.Run("when ticker is not found", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
		c.Request.Header.Set("Origin", "https://demoticker.org")
		store := &storage.MockStorage{}
		store.On("FindTickerByDomain", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
		mw := PrefetchTickerFromRequest(store)

		mw(c)

		s.Equal(http.StatusOK, w.Code)
		ticker, exists := c.Get("ticker")
		s.Nil(ticker)
		s.False(exists)
	})

	s.Run("when ticker is found", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
		c.Request.Header.Set("Origin", "https://demoticker.org")
		store := &storage.MockStorage{}
		store.On("FindTickerByDomain", mock.Anything).Return(storage.Ticker{}, nil)
		mw := PrefetchTickerFromRequest(store)

		mw(c)

		s.Equal(http.StatusOK, w.Code)
		ticker, exists := c.Get("ticker")
		s.NotNil(ticker)
		s.True(exists)
	})
}

func TestTickerTestSuite(t *testing.T) {
	suite.Run(t, new(TickerTestSuite))
}
