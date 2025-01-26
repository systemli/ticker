package response_cache

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/cache"
)

type ResponseCacheTestSuite struct {
	suite.Suite
}

func (s *ResponseCacheTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *ResponseCacheTestSuite) TestCreateKey() {
	s.Run("create cache key with origin", func() {
		c := gin.Context{
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/api/v1/settings", RawQuery: "origin=localhost"},
			},
		}

		key := CreateKey(&c)
		s.Equal("response:localhost:/api/v1/settings:origin=localhost", key)
	})

	s.Run("create cache key without origin", func() {
		c := gin.Context{
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/api/v1/settings"},
			},
		}

		key := CreateKey(&c)
		s.Equal("response:unknown:/api/v1/settings:", key)
	})
}

func (s *ResponseCacheTestSuite) TestCachePage() {
	s.Run("when cache is empty", func() {
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

		s.Equal(http.StatusOK, w.Code)
		s.Equal("pong", w.Body.String())

		count := 0
		inMemoryCache.Range(func(key, value interface{}) bool {
			count++
			return true
		})
		s.Equal(1, count)
	})

	s.Run("when cache is not empty", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Method: "GET",
			Header: http.Header{
				"Origin": []string{"http://localhost/"},
			},
			URL: &url.URL{Path: "/ping"},
		}

		inMemoryCache := cache.NewCache(time.Minute)
		defer inMemoryCache.Close()
		inMemoryCache.Set("response:http://localhost:/ping:", responseCache{
			Status: http.StatusOK,
			Header: http.Header{
				"DNT": []string{"1"},
			},
			Body: []byte("cached"),
		}, time.Minute)

		CachePage(inMemoryCache, time.Minute, func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})(c)

		s.Equal(http.StatusOK, w.Code)
		s.Equal("cached", w.Body.String())
	})
}

func TestResponseCacheTestSuite(t *testing.T) {
	suite.Run(t, new(ResponseCacheTestSuite))
}
