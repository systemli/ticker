package api

import (
	"errors"
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

type MediaTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *MediaTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.w)
	s.ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/media", nil)
	s.store = &storage.MockStorage{}
	s.cfg = config.LoadConfig("")
}

func (s *MediaTestSuite) TestGetMedia() {
	s.Run("when upload not found", func() {
		s.store.On("FindUploadByUUID", mock.Anything).Return(storage.Upload{}, errors.New("not found")).Once()
		h := s.handler()
		h.GetMedia(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.store.AssertExpectations(s.T())
	})

	s.Run("when upload found", func() {
		upload := storage.NewUpload("image/png", 1)
		uploadPath := s.T().TempDir()
		fullPath := upload.FullPath(uploadPath)
		s.NoError(os.MkdirAll(filepath.Dir(fullPath), 0750))
		s.NoError(os.WriteFile(fullPath, []byte("not really a png"), 0600))

		s.store.On("FindUploadByUUID", mock.Anything).Return(upload, nil).Once()
		s.store.On("UploadPath").Return(uploadPath).Once()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/media/"+upload.FileName(), nil)
		ctx.AddParam("fileName", upload.FileName())

		h := s.handler()
		h.GetMedia(ctx)

		s.Equal(http.StatusOK, w.Code)
		s.Equal("image/png", w.Header().Get("Content-Type"))
		s.Equal("nosniff", w.Header().Get("X-Content-Type-Options"))
		s.Equal("public, max-age=2592000, immutable", w.Header().Get("Cache-Control"))
		s.Contains(w.Header().Get("Content-Security-Policy"), "default-src 'none'")
		s.store.AssertExpectations(s.T())
	})
}

func (s *MediaTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestMediaTestSuite(t *testing.T) {
	suite.Run(t, new(MediaTestSuite))
}
