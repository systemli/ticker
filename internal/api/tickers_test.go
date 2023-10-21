package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-mastodon"
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
	s := &storage.MockStorage{}
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
	c.Set("me", storage.User{IsSuperAdmin: true})
	s := &storage.MockStorage{}
	s.On("FindTickers", mock.Anything).Return([]storage.Ticker{}, errors.New("storage error"))
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
	c.Set("me", storage.User{IsSuperAdmin: false, Tickers: []storage.Ticker{{ID: 2}}})
	s := &storage.MockStorage{}
	s.On("FindTickersByIDs", mock.Anything, mock.Anything).Return([]storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTickerUsersTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickerUsers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTickerUsers(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("FindUsersByTicker", mock.Anything, mock.Anything).Return([]storage.User{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetTickerUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostTickerFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", nil)
	s := &storage.MockStorage{}
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
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := `{"domain":"localhost","title":"title","description":"description"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
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
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := `{"domain":"localhost","title":"title","description":"description"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	body := `{"domain":"localhost","title":"title","description":"description"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	body := `{"domain":"localhost","title":"title","description":"description"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutTickerUsersNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	body := `{"users":[1,2,3]}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("FindUsersByIDs", mock.Anything).Return(nil, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerUsers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPutTickerUsersStorageError2(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	body := `{"users":[1,2,3]}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("FindUsersByIDs", mock.Anything).Return([]storage.User{}, nil)
	s.On("FindUsersByTicker", mock.Anything).Return([]storage.User{}, nil)
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
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
	c.Set("ticker", storage.Ticker{})
	body := `{"users":[1,2,3]}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("FindUsersByIDs", mock.Anything).Return([]storage.User{}, nil)
	s.On("FindUsersByTicker", mock.Anything).Return([]storage.User{}, nil)
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutTickerTelegramTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	body := `{"active":true,"channelName":"@channel_name"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	body := `{"active":true,"channelName":"@channel_name"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerTelegram(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTickerTelegramTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerTelegram(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTickerTelegramStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerTelegram(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTickerTelegram(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerTelegram(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutTickerMastodonTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerMastodon(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutTickerMastodonFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/mastodon", nil)
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerMastodon(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerMastodonConnectError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	body := `{"active":true,"server":"http://localhost","secret":"secret","token":"token","accessToken":"access_token"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/mastodon", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerMastodon(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerMastodonStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account := mastodon.Account{}
		json, _ := json.Marshal(account)
		w.Write(json)
	}))
	defer server.Close()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	body := fmt.Sprintf(`{"server":"%s","token":"token","secret":"secret","accessToken":"access_toklen"}`, server.URL)
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/mastodon", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerMastodon(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutTickerMastodon(t *testing.T) {
	w := httptest.NewRecorder()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account := mastodon.Account{}
		json, _ := json.Marshal(account)
		w.Write(json)
	}))
	defer server.Close()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	body := fmt.Sprintf(`{"server":"%s","token":"token","secret":"secret","accessToken":"access_toklen"}`, server.URL)
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/mastodon", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutTickerMastodon(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTickerMastodonTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerMastodon(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTickerMastodonStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerMastodon(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTickerMastodon(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("SaveTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerMastodon(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTickerTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("DeleteMessages", mock.Anything).Return(errors.New("storage error"))
	s.On("DeleteUploadsByTicker", mock.Anything).Return(errors.New("storage error"))
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
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("DeleteMessages", mock.Anything).Return(nil)
	s.On("DeleteUploadsByTicker", mock.Anything).Return(nil)
	s.On("DeleteTicker", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTickerUserTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	c.AddParam("userID", "1")
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	c.AddParam("userID", "1")
	s := &storage.MockStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	s.On("DeleteTickerUser", mock.Anything, mock.Anything).Return(errors.New("storage error"))
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
	c.Set("ticker", storage.Ticker{})
	c.AddParam("userID", "1")
	s := &storage.MockStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	s.On("DeleteTickerUser", mock.Anything, mock.Anything).Return(nil)
	s.On("FindUsersByTicker", mock.Anything).Return([]storage.User{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteTickerUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResetTickerUserTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
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
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("DeleteMessages", mock.Anything).Return(errors.New("storage error"))
	s.On("DeleteUploadsByTicker", mock.Anything).Return(errors.New("storage error"))
	s.On("SaveTicker", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.ResetTicker(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestResetTickerStorageError2(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("DeleteMessages", mock.Anything).Return(errors.New("storage error"))
	s.On("DeleteUploadsByTicker", mock.Anything).Return(errors.New("storage error"))
	s.On("SaveTicker", mock.Anything).Return(nil)
	s.On("DeleteTickerUsers", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.ResetTicker(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestResetTicker(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockStorage{}
	s.On("DeleteMessages", mock.Anything).Return(nil)
	s.On("DeleteUploadsByTicker", mock.Anything).Return(nil)
	s.On("SaveTicker", mock.Anything).Return(nil)
	s.On("DeleteTickerUsers", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.ResetTicker(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
