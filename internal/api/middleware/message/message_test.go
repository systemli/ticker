package message

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

type MessageTestSuite struct {
	suite.Suite
}

func (s *MessageTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *MessageTestSuite) TestMessage() {
	s.Run("when param is missing", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("ticker", storage.Ticker{})
		store := &storage.MockStorage{}
		mw := PrefetchMessage(store)

		mw(c)

		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("storage returns error", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("messageID", "1")
		c.Set("ticker", storage.Ticker{})
		store := &storage.MockStorage{}
		store.On("FindMessage", mock.Anything, mock.Anything, mock.Anything).Return(storage.Message{}, errors.New("storage error"))
		mw := PrefetchMessage(store)

		mw(c)

		s.Equal(http.StatusNotFound, w.Code)
	})

	s.Run("storage returns message", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("messageID", "1")
		c.Set("ticker", storage.Ticker{})
		store := &storage.MockStorage{}
		message := storage.Message{ID: 1}
		store.On("FindMessage", mock.Anything, mock.Anything, mock.Anything).Return(message, nil)
		mw := PrefetchMessage(store)

		mw(c)

		me, e := c.Get("message")
		s.True(e)
		s.Equal(message, me.(storage.Message))
	})
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}
