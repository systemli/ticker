package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetFeatures(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetFeatures(c)

	expected := `{"data":{"features":{"telegram_enabled":false}},"status":"success","error":{}}`
	assert.Equal(t, expected, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}
