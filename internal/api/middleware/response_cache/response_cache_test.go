package response_cache

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/cache"
)

func TestCreateKey(t *testing.T) {
	c := gin.Context{
		Request: &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: "/api/v1/settings", RawQuery: "origin=localhost"},
		},
	}

	key := CreateKey(&c)
	assert.Equal(t, "response:localhost::origin=localhost", key)

	c.Request.URL.RawQuery = ""

	key = CreateKey(&c)
	assert.Equal(t, "response:unknown::", key)
}

func TestCachePage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/ping", RawQuery: "origin=localhost"},
	}

	inMemoryCache := cache.NewCache(time.Minute)
	defer inMemoryCache.Close()
	CachePage(inMemoryCache, time.Minute, func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})(c)

	assert.Equal(t, http.StatusOK, w.Code)

	CachePage(inMemoryCache, time.Minute, func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
