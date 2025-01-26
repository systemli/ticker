package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type InitTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *InitTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.w)
	s.store = &storage.MockStorage{}
	s.store.On("GetRefreshIntervalSettings").Return(storage.DefaultRefreshIntervalSettings())
	s.cfg = config.LoadConfig("")
}

func (s *InitTestSuite) TestGetInit() {
	s.Run("if neither header nor query is set", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/init", nil)
		h := s.handler()
		h.GetInit(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Equal(`{"data":{"settings":{"refreshInterval":10000},"ticker":null},"status":"success","error":{}}`, s.w.Body.String())
		s.store.AssertNotCalled(s.T(), "FindTickerByOrigin", "", mock.Anything)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when database returns error", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/init?origin=https://demoticker.org", nil)
		s.store.On("FindTickerByOrigin", "https://demoticker.org", mock.Anything).Return(storage.Ticker{}, errors.New("storage error")).Once()
		s.store.On("GetInactiveSettings").Return(storage.DefaultInactiveSettings()).Once()
		h := s.handler()
		h.GetInit(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), `"ticker":null`)
		s.store.AssertCalled(s.T(), "FindTickerByOrigin", "https://demoticker.org", mock.Anything)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when database returns an inactive ticker", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/init?origin=https://demoticker.org", nil)
		ticker := storage.NewTicker()
		ticker.Active = false
		s.store.On("FindTickerByOrigin", "https://demoticker.org", mock.Anything).Return(ticker, nil).Once()
		s.store.On("GetInactiveSettings").Return(storage.DefaultInactiveSettings()).Once()
		h := s.handler()
		h.GetInit(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), `"ticker":null`)
		s.store.AssertCalled(s.T(), "FindTickerByOrigin", "https://demoticker.org", mock.Anything)
	})

	s.Run("when database returns an active ticker", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/init?origin=https://demoticker.org", nil)
		ticker := storage.NewTicker()
		ticker.Active = true
		s.store.On("FindTickerByOrigin", "https://demoticker.org", mock.Anything).Return(ticker, nil).Once()
		s.store.On("GetInactiveSettings").Return(storage.DefaultInactiveSettings()).Once()
		h := s.handler()
		h.GetInit(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertCalled(s.T(), "FindTickerByOrigin", "https://demoticker.org", mock.Anything)
	})
}

func (s *InitTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestInitTestSuite(t *testing.T) {
	suite.Run(t, new(InitTestSuite))
}
