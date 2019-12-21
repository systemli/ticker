package model_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
)

func TestNewUpload(t *testing.T) {
	u := initUploadTestData()

	assert.Equal(t, 0, u.ID)
	assert.Equal(t, "image/jpeg", u.ContentType)
	assert.Equal(t, 1, u.TickerID)
	assert.NotNil(t, u.Path)
	assert.NotNil(t, u.UUID)
	assert.NotNil(t, u.CreationDate)
}

func TestUpload_FileName(t *testing.T) {
	u := initUploadTestData()

	assert.Equal(t, fmt.Sprintf("%s.%s", u.UUID, u.Extension), u.FileName())
}

func TestUpload_FullPath(t *testing.T) {
	u := initUploadTestData()

	assert.Equal(t, fmt.Sprintf("uploads/%d/%d/%s", u.CreationDate.Year(), u.CreationDate.Month(), u.FileName()), u.FullPath())
}

func TestUpload_URL(t *testing.T) {
	u := initUploadTestData()

	assert.Equal(t, fmt.Sprintf("%s/media/%s", model.Config.UploadURL, u.FileName()), u.URL())
}

func TestNewUploadResponse(t *testing.T) {
	u := initUploadTestData()
	r := model.NewUploadResponse(u)

	assert.Equal(t, u.URL(), r.URL)
	assert.Equal(t, u.ContentType, r.ContentType)
}

func TestNewUploadsResponse(t *testing.T) {
	u := initUploadTestData()
	r := model.NewUploadsResponse([]*model.Upload{u})

	assert.Equal(t, 1, len(r))
}

func initUploadTestData() *model.Upload {
	model.Config = model.NewConfig()
	u := model.NewUpload("image.jpg", "image/jpeg", 1)

	return u
}
