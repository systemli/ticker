package api_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/api"
)

func TestUserMiddleware(t *testing.T) {
	setup()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		user := c.Query("user")
		if user != "" {
			uID, err := strconv.Atoi(user)
			if err != nil {
				t.Fail()
			}

			c.Set("id", float64(uID))
		}
	})
	router.Use(api.UserMiddleware())
	router.GET("/login", func(c *gin.Context) {
		c.String(200, "")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"data":{},"status":"error","error":{"code":1000,"message":"user identifier not found"}}`, strings.TrimSuffix(w.Body.String(), "\n"))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/login?user=2000", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"data":{},"status":"error","error":{"code":1000,"message":"user not found"}}`, strings.TrimSuffix(w.Body.String(), "\n"))
}
