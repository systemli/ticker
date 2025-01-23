package response_cache

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/cache"
)

// responseCache is a struct to cache the response
type responseCache struct {
	Status int
	Header http.Header
	Body   []byte
}

// cachedWriter is a wrapper around the gin.ResponseWriter
var _ gin.ResponseWriter = &cachedWriter{}

// cachedWriter is a wrapper around the gin.ResponseWriter
type cachedWriter struct {
	gin.ResponseWriter
	status  int
	written bool
	key     string
	expires time.Duration
	cache   *cache.Cache
}

// WriteHeader is a wrapper around the gin.ResponseWriter.WriteHeader
func (w *cachedWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// Status is a wrapper around the gin.ResponseWriter.Status
func (w *cachedWriter) Status() int {
	return w.ResponseWriter.Status()
}

// Written is a wrapper around the gin.ResponseWriter.Written
func (w *cachedWriter) Written() bool {
	return w.ResponseWriter.Written()
}

// Write is a wrapper around the gin.ResponseWriter.Write
// It will cache the response if the status code is below 300
func (w *cachedWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err == nil && w.Status() < 300 {
		value := responseCache{
			Status: w.Status(),
			Header: w.Header(),
			Body:   data,
		}
		w.cache.Set(w.key, value, w.expires)
	}

	return ret, err
}

// WriteString is a wrapper around the gin.ResponseWriter.WriteString
// It will cache the response if the status code is below 300
func (w *cachedWriter) WriteString(s string) (int, error) {
	ret, err := w.ResponseWriter.WriteString(s)
	if err == nil && w.Status() < 300 {
		value := responseCache{
			Status: w.Status(),
			Header: w.Header(),
			Body:   []byte(s),
		}
		w.cache.Set(w.key, value, w.expires)
	}

	return ret, err
}

func newCachedWriter(w gin.ResponseWriter, cache *cache.Cache, key string, expires time.Duration) *cachedWriter {
	return &cachedWriter{
		ResponseWriter: w,
		cache:          cache,
		key:            key,
		expires:        expires,
	}
}

// CachePage is a middleware to cache the response of a request
func CachePage(cache *cache.Cache, expires time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := CreateKey(c)
		if value, exists := cache.Get(key); exists {
			v := value.(responseCache)
			for k, values := range v.Header {
				for _, value := range values {
					c.Writer.Header().Set(k, value)
				}
			}
			c.Writer.WriteHeader(v.Status)
			_, _ = c.Writer.Write(v.Body)

			return
		} else {
			writer := newCachedWriter(c.Writer, cache, key, expires)
			c.Writer = writer
			handle(c)

			if c.IsAborted() {
				cache.Delete(key)
			}
		}
	}
}

func CreateKey(c *gin.Context) string {
	domain, err := helper.GetOrigin(c)
	if err != nil {
		domain = "unknown"
	}
	name := c.Request.URL.Path
	query := c.Request.URL.Query().Encode()

	return fmt.Sprintf("response:%s:%s:%s", domain, name, query)
}
