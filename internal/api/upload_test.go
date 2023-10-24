package api

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestPostUploadForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostUploadMultipartError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", nil)
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUploadMissingTickerValue(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.CreateFormField("field")
	_ = writer.Close()
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUploadTickerValueWrong(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.CreateFormField("ticker")
	_ = writer.Close()
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUploadTickerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("ticker", "1")
	_ = writer.Close()
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	s := &storage.MockStorage{}
	s.On("FindTickerByUserAndID", mock.Anything, mock.Anything).Return(storage.Ticker{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUploadMissingFiles(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("ticker", "1")
	_ = writer.Close()
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	s := &storage.MockStorage{}
	s.On("FindTickerByUserAndID", mock.Anything, mock.Anything).Return(storage.Ticker{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUpload(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("ticker", "1")
	path := "../../testdata/gopher.jpg"
	part, err := writer.CreateFormFile("files", filepath.Base(path))
	if err != nil {
		t.Error(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	part.Write(b)
	_ = writer.Close()
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	s := &storage.MockStorage{}
	s.On("FindTickerByUserAndID", mock.Anything, mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveUpload", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUploadGIF(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("ticker", "1")
	path := "../../testdata/gopher-dance.gif"
	part, err := writer.CreateFormFile("files", filepath.Base(path))
	if err != nil {
		t.Error(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	part.Write(b)
	_ = writer.Close()
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	s := &storage.MockStorage{}
	s.On("FindTickerByUserAndID", mock.Anything, mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveUpload", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUploadTooMuchFiles(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("ticker", "1")
	path := "../../testdata/gopher.jpg"
	part1, _ := writer.CreateFormFile("files", filepath.Base(path))
	part2, _ := writer.CreateFormFile("files", filepath.Base(path))
	part3, _ := writer.CreateFormFile("files", filepath.Base(path))
	part4, _ := writer.CreateFormFile("files", filepath.Base(path))
	b, _ := os.ReadFile(path)

	part1.Write(b)
	part2.Write(b)
	part3.Write(b)
	part4.Write(b)
	_ = writer.Close()
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	s := &storage.MockStorage{}
	s.On("FindTickerByUserAndID", mock.Anything, mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveUpload", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUploadForbiddenFileType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("ticker", "1")
	path := "./api.go"
	part, err := writer.CreateFormFile("files", filepath.Base(path))
	if err != nil {
		t.Error(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	part.Write(b)
	_ = writer.Close()
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	s := &storage.MockStorage{}
	s.On("FindTickerByUserAndID", mock.Anything, mock.Anything).Return(storage.Ticker{}, nil)
	s.On("SaveUpload", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUpload(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
