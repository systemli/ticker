package api_test

import (
	"testing"


	"github.com/gin-gonic/gin"
	"git.codecoop.org/systemli/ticker/internal/api"
	"net/http/httptest"
	"net/http"
	"github.com/stretchr/testify/assert"
)

func TestUserMiddleware(t *testing.T) {
	setup()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		user := c.Query("user")
		if user != "" {
			c.Set("userID", user)
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
	assert.Equal(t, `{"data":{},"status":"error","error":{"code":1000,"message":"user identifier not found"}}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/login?user=2000", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"data":{},"status":"error","error":{"code":1000,"message":"user not found"}}`, w.Body.String())
}
