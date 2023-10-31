package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
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
	s.ctx.Request = httptest.NewRequest(http.MethodGet, "/media", nil)
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
		upload := storage.NewUpload("image.jpg", "image/jpeg", 1)
		s.store.On("FindUploadByUUID", mock.Anything).Return(upload, nil).Once()
		s.store.On("UploadPath").Return("./uploads").Once()

		h := s.handler()
		h.GetMedia(s.ctx)

		s.Equal(http.StatusNotFound, s.w.Code)
		s.NotEmpty(s.w.Header().Get("Cache-Control"))
		s.NotEmpty(s.w.Header().Get("Expires"))
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
