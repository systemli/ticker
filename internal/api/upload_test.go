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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type UploadTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *UploadTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *UploadTestSuite) Run(name string, subtest func()) {
	s.T().Run(name, func(t *testing.T) {
		s.w = httptest.NewRecorder()
		s.ctx, _ = gin.CreateTestContext(s.w)
		s.store = &storage.MockStorage{}
		s.cfg = config.LoadConfig("")

		subtest()
	})
}

func (s *UploadTestSuite) TestPostUpload() {
	s.Run("when no user is set", func() {
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when form is invalid", func() {
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", nil)
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when ticker is missing", func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.CreateFormField("field")
		_ = writer.Close()
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when ticker value is invalid", func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.CreateFormField("ticker")
		_ = writer.Close()
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when ticker is not found", func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("ticker", "1")
		_ = writer.Close()
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickerByUserAndID", mock.Anything, 1).Return(storage.Ticker{}, errors.New("not found")).Once()
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when form files are missing", func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("ticker", "1")
		_ = writer.Close()
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickerByUserAndID", mock.Anything, 1).Return(storage.Ticker{}, nil).Once()
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when too much files are uploaded", func() {
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
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickerByUserAndID", mock.Anything, 1).Return(storage.Ticker{}, nil).Once()
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when file type is not allowed", func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("ticker", "1")
		path := "./api.go"
		part, _ := writer.CreateFormFile("files", filepath.Base(path))
		b, _ := os.ReadFile(path)
		part.Write(b)
		_ = writer.Close()
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickerByUserAndID", mock.Anything, 1).Return(storage.Ticker{}, nil).Once()
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when file is gif", func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("ticker", "1")
		path := "../../testdata/gopher-dance.gif"
		part, _ := writer.CreateFormFile("files", filepath.Base(path))
		b, _ := os.ReadFile(path)
		part.Write(b)
		_ = writer.Close()
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickerByUserAndID", mock.Anything, 1).Return(storage.Ticker{}, nil).Once()
		s.store.On("SaveUpload", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when save returns an error", func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("ticker", "1")
		path := "../../testdata/gopher.jpg"
		part, _ := writer.CreateFormFile("files", filepath.Base(path))
		b, _ := os.ReadFile(path)
		part.Write(b)
		_ = writer.Close()
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickerByUserAndID", mock.Anything, 1).Return(storage.Ticker{}, nil).Once()
		s.store.On("SaveUpload", mock.Anything).Return(errors.New("save error")).Once()
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusBadRequest, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when save is successful", func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("ticker", "1")
		path := "../../testdata/gopher.jpg"
		part, _ := writer.CreateFormFile("files", filepath.Base(path))
		b, _ := os.ReadFile(path)
		part.Write(b)
		_ = writer.Close()
		s.ctx.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
		s.ctx.Request.Header.Add("Content-Type", writer.FormDataContentType())
		s.ctx.Set("me", storage.User{IsSuperAdmin: true})
		s.store.On("FindTickerByUserAndID", mock.Anything, 1).Return(storage.Ticker{}, nil).Once()
		s.store.On("SaveUpload", mock.Anything).Return(nil).Once()
		h := s.handler()
		h.PostUpload(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.store.AssertExpectations(s.T())
	})
}

func (s *UploadTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestUploadTestSuite(t *testing.T) {
	suite.Run(t, new(UploadTestSuite))
}
