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

func TestGetMessagesForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetMessagesWithoutParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMessagesWithoutPermissions(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetMessagesTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindMessagesByTicker", mock.Anything, mock.Anything).Return([]storage.Message{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetMessages(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindMessagesByTicker", mock.Anything, mock.Anything).Return([]storage.Message{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessages(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetMessageForbidden(t *testing.T) {
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

func TestGetMessageWithoutTickerParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetMessageWithoutPermissions(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessage(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetMessageWithoutWithoutMessageParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetMessageStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.AddParam("messageID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindMessage", mock.Anything, mock.Anything).Return(storage.Message{}, errors.New("storage error"))
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.AddParam("messageID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindMessage", mock.Anything, mock.Anything).Return(storage.Message{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetMessage(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostMessageForbidden(t *testing.T) {
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostMessageMissingTickerParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	json := `{"text":"text"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostMessageWithoutPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false})
	json := `{"text":"text"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPostMessageTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	json := `{"text":"text"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostMessageUploadsNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	json := `{"text":"text","attachments":[1]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	json := `{"text":"text","attachments":[1]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	json := `{"text":"text","attachments":[1]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
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

func TestDeleteMessageForbidden(t *testing.T) {
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

func TestDeleteMessageMissingTickerParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteMessage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteMessageWithoutPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteMessage(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteMessageTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteMessage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteMessageMissingMessageParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteMessage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteMessageMessageNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.AddParam("messageID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindMessage", mock.Anything, mock.Anything).Return(storage.Message{}, errors.New("not found"))
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.AddParam("messageID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindMessage", mock.Anything, mock.Anything).Return(storage.Message{}, nil)
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
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.AddParam("messageID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindMessage", mock.Anything, mock.Anything).Return(storage.Message{}, nil)
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
