package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetMedia(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/media", nil)

	upload := storage.NewUpload("image.jpg", "image/jpeg", 1)
	s := &storage.MockStorage{}
	s.On("FindUploadByUUID", mock.Anything).Return(upload, nil)
	s.On("UploadPath").Return("./uploads")

	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetMedia(c)

	assert.NotEmpty(t, w.Header().Get("Cache-Control"))
	assert.NotEmpty(t, w.Header().Get("Expires"))
}

func TestGetMediaNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	s := &storage.MockStorage{}
	s.On("FindUploadByUUID", mock.Anything).Return(storage.Upload{}, errors.New("not found"))

	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}
	h.GetMedia(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
