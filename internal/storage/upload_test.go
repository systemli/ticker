package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var upload = NewUpload("image/jpeg", 1)

func TestUploadFilename(t *testing.T) {
	fileName := upload.FileName()

	assert.Contains(t, fileName, ".jpg")
}

func TestUploadFullPath(t *testing.T) {
	fullPath := upload.FullPath("/uploads")
	assert.Contains(t, fullPath, "/uploads")
}

func TestUploadURL(t *testing.T) {
	assert.Equal(t, "/api/media/"+upload.FileName(), upload.URL())
}

func TestExtensionForContentType(t *testing.T) {
	ext, ok := ExtensionForContentType("image/png")
	assert.True(t, ok)
	assert.Equal(t, "png", ext)

	ext, ok = ExtensionForContentType("text/html")
	assert.False(t, ok)
	assert.Empty(t, ext)
}

func TestNewUploadIgnoresClientExtension(t *testing.T) {
	// A GIF is stored unmodified, so a filename like "evil.html" used to decide
	// how the file is served afterwards.
	u := NewUpload("image/gif", 1)

	assert.Equal(t, "gif", u.Extension)
	assert.Equal(t, u.UUID+".gif", u.FileName())
}
