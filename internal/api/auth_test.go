package api_test

import (
	"testing"

	"github.com/systemli/ticker/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
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

			c.Set("userID", float64(uID))
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
