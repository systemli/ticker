package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNeedAdminMissingUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mw := NeedAdmin()

	mw(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNeedAdminNonAdmin(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{})
	mw := NeedAdmin()

	mw(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
