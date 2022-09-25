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
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetTickersForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTickersStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("FindTickers").Return([]storage.Ticker{}, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTickers(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false, Tickers: []int{1}})
	s := &storage.MockTickerStorage{}
	s.On("FindTickersByIDs", mock.Anything).Return([]storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTickerForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTicker(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTickerMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTicker(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTickerMissingPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false, Tickers: []int{1}})
	c.AddParam("tickerID", "2")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTicker(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetTickerStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false, Tickers: []int{1}})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTicker(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTicker(t *testing.T) {
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

	h.GetTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTickerUsersForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickerUsers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTickerUsersMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickerUsers(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTickerUsersTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickerUsers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTickerUsersWrongPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: false})
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickerUsers(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetTickerUsers(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindUsersByTicker", mock.Anything).Return([]storage.User{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickerUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostTickerForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostTicker(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPostTickerFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", nil)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostTicker(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostTickerStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	body := `{"domain":"localhost","title":"title","description":"description"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostTicker(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostTicker(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	body := `{"domain":"localhost","title":"title","description":"description"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutTickerForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTicker(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutTickerMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTicker(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerWrongPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: false})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTicker(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPutTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTicker(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutTickerFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTicker(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	body := `{"domain":"localhost","title":"title","description":"description"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTicker(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTicker(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	body := `{"domain":"localhost","title":"title","description":"description"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutTickerUsersForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerUsers(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPutTickerUsersMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerUsers(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerUsersNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerUsers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutTickerUsersFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerUsers(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerUsersStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	body := `{"users":[1,2,3]}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("AddUsersToTicker", mock.Anything, mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerUsers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPutTickerUsers(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("tickerID", "1")
	c.Set("user", storage.User{IsSuperAdmin: true})
	body := `{"users":[1,2,3]}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("AddUsersToTicker", mock.Anything, mock.Anything).Return(nil)
	s.On("FindUsersByTicker", mock.Anything).Return([]storage.User{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutTickerTwitterForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTwitter(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutTickerTwitterMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTwitter(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerTwitterWrongPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTwitter(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPutTickerTwitterTickerNotFound(t *testing.T) {
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

	h.PutTickerTwitter(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutTickerTwitterFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/twitter", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTwitter(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerTwitterDisconnect(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	body := `{"disconnect":true}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/twitter", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTwitter(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutTickerTwitterConnect(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	body := `{"active":true,"token":"token","secret":"secret"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/twitter", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTwitter(c)
}

func TestPutTickerTwitterStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	body := `{"disconnect":true}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/twitter", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTwitter(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerTelegramForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTelegram(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutTickerTelegramMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTelegram(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerTelegramWrongPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: false})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTelegram(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPutTickerTelegramTickerNotFound(t *testing.T) {
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

	h.PutTickerTelegram(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutTickerTelegramFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTelegram(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerTelegramStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	body := `{"active":true,"channel_name":"@channel_name"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTelegram(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerTelegram(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	body := `{"active":true,"channel_name":"@channel_name"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTelegram(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTickerForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTicker(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteTickerMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTicker(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTickerTickerNotFound(t *testing.T) {
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

	h.DeleteTicker(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTickerStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("DeleteMessages", mock.Anything).Return(errors.New("storage error"))
	s.On("DeleteTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTicker(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTicker(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("DeleteMessages", mock.Anything).Return(nil)
	s.On("DeleteTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTickerUserForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerUser(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteTickerMissingTickerParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTickerUserTickerNotFound(t *testing.T) {
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

	h.DeleteTickerUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTickerUserMissingUserParam(t *testing.T) {
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

	h.DeleteTickerUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTickerUserUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.AddParam("userID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTickerUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.AddParam("userID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	s.On("RemoveTickerFromUser", mock.Anything, mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteTickerUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	c.AddParam("userID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	s.On("RemoveTickerFromUser", mock.Anything, mock.Anything).Return(nil)
	s.On("FindUsersByTicker", mock.Anything).Return([]storage.User{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResetTickerForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.ResetTicker(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestResetTickerMissingTickerParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.ResetTicker(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetTickerUserTickerNotFound(t *testing.T) {
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

	h.ResetTicker(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestResetTickerStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("DeleteMessages", mock.Anything).Return(errors.New("storage error"))
	s.On("DeleteUploadsByTicker", mock.Anything).Return(errors.New("storage error"))
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.ResetTicker(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestResetTicker(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	c.AddParam("tickerID", "1")
	s := &storage.MockTickerStorage{}
	s.On("FindTickerByID", mock.Anything).Return(storage.Ticker{}, nil)
	s.On("DeleteMessages", mock.Anything).Return(nil)
	s.On("DeleteUploadsByTicker", mock.Anything).Return(nil)
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.ResetTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
