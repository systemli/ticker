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
	s := &storage.MockStorage{}

	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetFeatures(c)

	expected := `{"data":{"features":{"telegramEnabled":false}},"status":"success","error":{}}`
	assert.Equal(t, expected, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}
