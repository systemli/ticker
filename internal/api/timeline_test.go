package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type TimelineTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *TimelineTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *TimelineTestSuite) Run(name string, subtest func()) {
	s.T().Run(name, func(t *testing.T) {
		s.w = httptest.NewRecorder()
		s.ctx, _ = gin.CreateTestContext(s.w)
		s.store = &storage.MockStorage{}
		s.cfg = config.LoadConfig("")

		subtest()
	})
}

func (s *TimelineTestSuite) TestGetTimeline() {
	s.Run("when ticker is missing", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
		h := s.handler()
		h.GetTimeline(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), response.TickerNotFound)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns an error", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
		s.ctx.Set("ticker", storage.Ticker{Active: true})
		s.store.On("FindMessagesByTickerAndPagination", mock.Anything, mock.Anything, mock.Anything).Return([]storage.Message{}, errors.New("storage error")).Once()
		h := s.handler()
		h.GetTimeline(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), response.MessageFetchError)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns messages", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/timeline", nil)
		s.ctx.Set("ticker", storage.Ticker{Active: true})
		s.store.On("FindMessagesByTickerAndPagination", mock.Anything, mock.Anything, mock.Anything).Return([]storage.Message{}, nil).Once()
		h := s.handler()
		h.GetTimeline(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TimelineTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestTimelineTestSuite(t *testing.T) {
	suite.Run(t, new(TimelineTestSuite))
}
