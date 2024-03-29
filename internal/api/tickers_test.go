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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type TickerTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *TickerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *TickerTestSuite) Run(name string, subtest func()) {
	s.T().Run(name, func(t *testing.T) {
		s.w = httptest.NewRecorder()
		s.ctx, _ = gin.CreateTestContext(s.w)
		s.store = &storage.MockStorage{}
		s.cfg = config.LoadConfig("")

		subtest()
	})
}

func (s *TickerTestSuite) TestGetTickers() {
	s.Run("when not authorized", func() {
		h := s.handler()
		h.GetTickers(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns an error", func() {
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickersByUser", mock.Anything, mock.Anything).Return([]storage.Ticker{}, errors.New("storage error")).Once()
		h := s.handler()
		h.GetTickers(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns tickers", func() {
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickersByUser", mock.Anything, mock.Anything).Return([]storage.Ticker{}, nil).Once()
		h := s.handler()
		h.GetTickers(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestGetTicker() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.GetTicker(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when ticker found", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		h := s.handler()
		h.GetTicker(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestGetTickerUsers() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.GetTickerUsers(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when ticker found", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("FindUsersByTicker", mock.Anything, mock.Anything).Return([]storage.User{}, nil).Once()
		h := s.handler()
		h.GetTickerUsers(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestPostTicker() {
	s.Run("when body is invalid", func() {
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", strings.NewReader(`broken_json`))
		h := s.handler()
		h.PostTicker(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		body := `{"domain":"localhost","title":"title","description":"description"}`
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PostTicker(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		body := `{"domain":"localhost","title":"title","description":"description"}`
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/tickers", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PostTicker(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestPutTicker() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.PutTicker(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when body is invalid", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", nil)
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutTicker(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		body := `{"domain":"localhost","title":"title","description":"description"}`
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutTicker(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		body := `{"domain":"localhost","title":"title","description":"description"}`
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PutTicker(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestPutTickerUsers() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.PutTickerUsers(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when body is invalid", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", nil)
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutTickerUsers(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		body := `{"users":[{"id":1},{"id":2},{"id":3}]}`
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutTickerUsers(s.ctx)

		s.Equal(http.StatusInternalServerError, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		body := `{"users":[{"id":1},{"id":2},{"id":3}]}`
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/user", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PutTickerUsers(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestPutTickerTelegram() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.PutTickerTelegram(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when body is invalid", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", nil)
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutTickerTelegram(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		body := `{"active":true,"channelName":"@channel_name"}`
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutTickerTelegram(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		body := `{"active":true,"channelName":"@channel_name"}`
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/telegram", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PutTickerTelegram(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestDeleteTickerTelegram() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.DeleteTickerTelegram(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("SaveTicker", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.DeleteTickerTelegram(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.DeleteTickerTelegram(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestPutTickerMastodon() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.PutTickerMastodon(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when body is invalid", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/mastodon", nil)
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutTickerMastodon(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when mastodon server is not reachable", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		body := `{"active":true,"server":"http://localhost","secret":"secret","token":"token","accessToken":"access_token"}`
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/mastodon", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		h := s.handler()
		h.PutTickerMastodon(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			account := mastodon.Account{}
			json, _ := json.Marshal(account)
			w.Write(json)
		}))
		defer server.Close()
		body := fmt.Sprintf(`{"server":"%s","token":"token","secret":"secret","accessToken":"access_toklen"}`, server.URL)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/mastodon", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.PutTickerMastodon(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			account := mastodon.Account{}
			json, _ := json.Marshal(account)
			w.Write(json)
		}))
		defer server.Close()
		body := fmt.Sprintf(`{"server":"%s","token":"token","secret":"secret","accessToken":"access_toklen"}`, server.URL)
		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/tickers/1/mastodon", strings.NewReader(body))
		s.ctx.Request.Header.Add("Content-Type", "application/json")
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PutTickerMastodon(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestDeleteTickerMastodon() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.DeleteTickerMastodon(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("SaveTicker", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.DeleteTickerMastodon(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.DeleteTickerMastodon(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestDeleteTicker() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.DeleteTicker(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("DeleteMessages", mock.Anything).Return(errors.New("storage error"))
		s.store.On("DeleteUploadsByTicker", mock.Anything).Return(errors.New("storage error"))
		s.store.On("DeleteTicker", mock.Anything).Return(errors.New("storage error"))
		h := s.handler()
		h.DeleteTicker(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("DeleteMessages", mock.Anything).Return(nil)
		s.store.On("DeleteUploadsByTicker", mock.Anything).Return(nil)
		s.store.On("DeleteTicker", mock.Anything).Return(nil)
		h := s.handler()
		h.DeleteTicker(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestDeleteTickerUser() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.DeleteTickerUser(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when user param is missing", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		h := s.handler()
		h.DeleteTickerUser(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when user not found", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.ctx.AddParam("userID", "1")
		s.store.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("not found")).Once()
		h := s.handler()
		h.DeleteTickerUser(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.ctx.AddParam("userID", "1")
		s.store.On("FindUserByID", mock.Anything).Return(storage.User{}, nil).Once()
		s.store.On("DeleteTickerUser", mock.Anything, mock.Anything).Return(errors.New("storage error"))
		h := s.handler()
		h.DeleteTickerUser(s.ctx)

		s.Equal(http.StatusInternalServerError, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.ctx.AddParam("userID", "1")
		s.store.On("FindUserByID", mock.Anything).Return(storage.User{}, nil).Once()
		s.store.On("DeleteTickerUser", mock.Anything, mock.Anything).Return(nil)
		h := s.handler()
		h.DeleteTickerUser(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) TestResetTicker() {
	s.Run("when ticker not found", func() {
		h := s.handler()
		h.ResetTicker(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns error", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("DeleteMessages", mock.Anything).Return(errors.New("storage error")).Once()
		s.store.On("DeleteUploadsByTicker", mock.Anything).Return(errors.New("storage error")).Once()
		s.store.On("SaveTicker", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.ResetTicker(s.ctx)

		s.Equal(http.StatusInternalServerError, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when deleting users fails", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("DeleteMessages", mock.Anything).Return(nil).Once()
		s.store.On("DeleteUploadsByTicker", mock.Anything).Return(nil).Once()
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		s.store.On("DeleteTickerUsers", mock.Anything).Return(errors.New("storage error")).Once()
		h := s.handler()
		h.ResetTicker(s.ctx)

		s.Equal(http.StatusInternalServerError, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when storage returns ticker", func() {
		s.ctx.Set("ticker", storage.Ticker{})
		s.store.On("DeleteMessages", mock.Anything).Return(nil).Once()
		s.store.On("DeleteUploadsByTicker", mock.Anything).Return(nil).Once()
		s.store.On("SaveTicker", mock.Anything).Return(nil).Once()
		s.store.On("DeleteTickerUsers", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.ResetTicker(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *TickerTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestTickerTestSuite(t *testing.T) {
	suite.Run(t, new(TickerTestSuite))
}
