package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/api/realtime"
	"github.com/systemli/ticker/internal/cache"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type MessagesTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	cache *cache.Cache
	suite.Suite
}

func (s *MessagesTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *MessagesTestSuite) Run(name string, subtest func()) {
	s.T().Run(name, func(t *testing.T) {
		s.w = httptest.NewRecorder()
		s.ctx, _ = gin.CreateTestContext(s.w)
		s.store = &storage.MockStorage{}
		s.cfg = config.LoadConfig("")
		s.cache = cache.NewCache(time.Minute)

		subtest()
	})
}

func (s *MessagesTestSuite) TestGetMessages() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.GetMessages(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("when database returns error", func() {
		ticker := storage.Ticker{ID: 1}
		s.ctx.Set("ticker", ticker)
		s.store.On("FindMessagesByTickerAndPagination", ticker, mock.Anything, mock.Anything).Return([]storage.Message{}, errors.New("storage error")).Once()
		h := s.handler()
		h.GetMessages(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when database returns messages", func() {
		ticker := storage.Ticker{ID: 1}
		s.ctx.Set("ticker", ticker)
		s.store.On("FindMessagesByTickerAndPagination", ticker, mock.Anything, mock.Anything).Return([]storage.Message{}, nil).Once()
		h := s.handler()
		h.GetMessages(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *MessagesTestSuite) TestGetMessage() {
	s.Run("when message not found", func() {
		h := s.handler()
		h.GetMessage(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("when message found", func() {
		message := storage.Message{ID: 1}
		s.ctx.Set("message", message)
		h := s.handler()
		h.GetMessage(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *MessagesTestSuite) TestPostMessage() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.PostMessage(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("when form is invalid", func() {
		ticker := storage.Ticker{ID: 1}
		s.ctx.Set("ticker", ticker)
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
		h := s.handler()
		h.PostMessage(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("when upload not found", func() {
		ticker := storage.Ticker{ID: 1}
		s.ctx.Set("ticker", ticker)
		json := `{"text":"text","attachments":[1]}`
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("FindUploadsByIDs", []int{1}).Return([]storage.Upload{}, errors.New("storage error")).Once()
		h := s.handler()
		h.PostMessage(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("when database returns error", func() {
		ticker := storage.Ticker{ID: 1}
		s.ctx.Set("ticker", ticker)
		json := `{"text":"text","attachments":[1]}`
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.ctx.AddParam("tickerID", "1")
		s.store.On("FindUploadsByIDs", []int{1}).Return([]storage.Upload{}, nil).Once()
		s.store.On("SaveMessage", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PostMessage(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("when database returns message", func() {
		ticker := storage.Ticker{ID: 1}
		s.ctx.Set("ticker", ticker)
		json := `{"text":"text","attachments":[1]}`
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.ctx.AddParam("tickerID", "1")
		s.store.On("FindUploadsByIDs", []int{1}).Return([]storage.Upload{}, nil).Once()
		s.store.On("SaveMessage", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PostMessage(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})
}

func (s *MessagesTestSuite) TestDeleteMessage() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.DeleteMessage(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("when message not found", func() {
		ticker := storage.Ticker{ID: 1}
		s.ctx.Set("ticker", ticker)
		h := s.handler()
		h.DeleteMessage(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("when database returns error", func() {
		ticker := storage.Ticker{ID: 1}
		message := storage.Message{ID: 1}
		s.ctx.Set("ticker", ticker)
		s.ctx.Set("message", message)
		s.store.On("DeleteMessage", message).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.DeleteMessage(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.True(s.store.AssertExpectations(s.T()))
	})

	s.Run("happy path", func() {
		ticker := storage.Ticker{ID: 1, Domain: "localhost"}
		message := storage.Message{ID: 1}
		s.cache.Set("response:localhost:/v1/timeline", true, time.Minute)
		s.ctx.Set("ticker", ticker)
		s.ctx.Set("message", message)
		s.store.On("DeleteMessage", message).Return(nil).Once()
		h := s.handler()
		h.DeleteMessage(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Nil(s.cache.Get("response:localhost:/v1/timeline"))
		s.True(s.store.AssertExpectations(s.T()))
	})
}

func (s *MessagesTestSuite) handler() handler {
	return handler{
		storage:  s.store,
		config:   s.cfg,
		cache:    s.cache,
		realtime: realtime.New(),
	}
}

func TestMessagesTestSuite(t *testing.T) {
	suite.Run(t, new(MessagesTestSuite))
}
