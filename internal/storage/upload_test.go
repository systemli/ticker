package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var upload = NewUpload("image.jpg", "image/jpeg", 1)

func TestUploadFilename(t *testing.T) {
	fileName := upload.FileName()

	assert.Contains(t, fileName, ".jpg")
}

func TestUploadFullPath(t *testing.T) {
	fullPath := upload.FullPath("/uploads")
	assert.Contains(t, fullPath, "/uploads")
}

func TestUploadURL(t *testing.T) {
	url := upload.URL("/uploads")
	assert.Contains(t, url, "/uploads/media")
}
