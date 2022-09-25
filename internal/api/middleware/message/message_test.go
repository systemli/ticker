package message

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

func TestPrefetchMessageParamMissing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	mw := PrefetchMessage(s)

	mw(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPrefetchMessageStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("messageID", "1")
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	s.On("FindMessage", mock.Anything, mock.Anything).Return(storage.Message{}, errors.New("storage error"))
	mw := PrefetchMessage(s)

	mw(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPrefetchMessage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("messageID", "1")
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	message := storage.Message{ID: 1}
	s.On("FindMessage", mock.Anything, mock.Anything).Return(message, nil)
	mw := PrefetchMessage(s)

	mw(c)

	me, e := c.Get("message")
	assert.True(t, e)
	assert.Equal(t, message, me.(storage.Message))

}
