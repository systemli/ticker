package response

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

var (
	u = storage.NewUpload("image.jpg", "image/jpg", 1)
	c = config.Config{
		Upload: config.Upload{URL: "http://localhost:8080"},
	}
)

func TestUploadResponse(t *testing.T) {
	response := UploadResponse(u, c)

	assert.Equal(t, fmt.Sprintf("%s/media/%s", c.Upload.URL, u.FileName()), response.URL)
}

func TestUploadsResponse(t *testing.T) {
	response := UploadsResponse([]storage.Upload{u}, c)

	assert.Equal(t, 1, len(response))
}
