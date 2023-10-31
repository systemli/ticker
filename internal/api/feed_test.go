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

type FeedTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *FeedTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.w)
	s.store = &storage.MockStorage{}
	s.cfg = config.LoadConfig("")
}

func (s *FeedTestSuite) TestGetFeed() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.GetFeed(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), response.TickerNotFound)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when fetching messages fails", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("FindMessagesByTickerAndPagination", mock.Anything, mock.Anything).Return([]storage.Message{}, errors.New("storage error")).Once()

		h := s.handler()
		h.GetFeed(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), response.MessageFetchError)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when fetching messages succeeds", func() {
		ticker := storage.Ticker{
			ID:    1,
			Title: "Title",
			Information: storage.TickerInformation{
				URL:    "https://demoticker.org",
				Author: "Author",
				Email:  "author@demoticker.org",
			},
		}
		s.ctx.Set("ticker", ticker)
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/feed?format=atom", nil)
		message := storage.Message{
			TickerID: ticker.ID,
			Text:     "Text",
		}
		s.store.On("FindMessagesByTickerAndPagination", mock.Anything, mock.Anything).Return([]storage.Message{message}, nil).Once()

		h := s.handler()
		h.GetFeed(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *FeedTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestFeedTestSuite(t *testing.T) {
	suite.Run(t, new(FeedTestSuite))
}
