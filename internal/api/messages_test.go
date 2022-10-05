package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetMessagesTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetMessagesStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	s.On("FindMessagesByTicker", mock.Anything, mock.Anything).Return([]storage.Message{}, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetMessagesEmptyResult(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	s.On("FindMessagesByTicker", mock.Anything, mock.Anything).Return([]storage.Message{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetMessages(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	s.On("FindMessagesByTicker", mock.Anything, mock.Anything).Return([]storage.Message{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetMessageNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetMessage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("message", storage.Message{})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostMessageTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostMessageFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostMessageUploadsNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	json := `{"text":"text","attachments":[1]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindUploadsByIDs", mock.Anything).Return([]storage.Upload{}, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostMessageStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	json := `{"text":"text","attachments":[1]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindUploadsByIDs", mock.Anything).Return([]storage.Upload{}, nil)
	s.On("SaveMessage", mock.Anything).Return(errors.New("storage error"))
	b := &bridge.MockBridge{}
	b.On("Send", mock.Anything, mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
		bridges: bridge.Bridges{"mock": b},
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostMessage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	json := `{"text":"text","attachments":[1]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindUploadsByIDs", mock.Anything).Return([]storage.Upload{}, nil)
	s.On("SaveMessage", mock.Anything).Return(nil)
	b := &bridge.MockBridge{}
	b.On("Send", mock.Anything, mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
		bridges: bridge.Bridges{"mock": b},
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteMessageTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteMessageMessageNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteMessageStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	c.Set("message", storage.Message{})
	s := &storage.MockTickerStorage{}
	s.On("DeleteMessage", mock.Anything).Return(errors.New("storage error"))
	b := &bridge.MockBridge{}
	b.On("Delete", mock.Anything, mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
		bridges: bridge.Bridges{"mock": b},
	}

	h.DeleteMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteMessage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	c.Set("message", storage.Message{})
	s := &storage.MockTickerStorage{}
	s.On("DeleteMessage", mock.Anything).Return(nil)
	b := &bridge.MockBridge{}
	b.On("Delete", mock.Anything, mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
		bridges: bridge.Bridges{"mock": b},
	}

	h.DeleteMessage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
