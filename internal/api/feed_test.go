package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetFeedTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetFeed(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), response.TickerNotFound)
}

func TestGetFeedMessageFetchError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("ticker", storage.Ticker{})
	s := &storage.MockTickerStorage{}
	s.On("FindMessagesByTickerAndPagination", mock.Anything, mock.Anything).Return([]storage.Message{}, errors.New("storage error"))

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetFeed(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), response.MessageFetchError)
}

func TestGetFeed(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ticker := storage.Ticker{
		ID:    1,
		Title: "Title",
		Information: storage.Information{
			URL:    "https://demoticker.org",
			Author: "Author",
			Email:  "author@demoticker.org",
		},
	}
	c.Set("ticker", ticker)
	s := &storage.MockTickerStorage{}
	message := storage.Message{
		Ticker:       ticker.ID,
		CreationDate: time.Now(),
		Text:         "Text",
	}
	s.On("FindMessagesByTickerAndPagination", mock.Anything, mock.Anything).Return([]storage.Message{message}, nil)

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetFeed(c)
}
