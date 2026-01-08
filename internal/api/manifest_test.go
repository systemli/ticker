package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type ManifestTestSuite struct {
	w     *httptest.ResponseRecorder
	ctx   *gin.Context
	store *storage.MockStorage
	cfg   config.Config
	suite.Suite
}

func (s *ManifestTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.w)
	s.store = &storage.MockStorage{}
	s.cfg = config.LoadConfig("")
}

func (s *ManifestTestSuite) TestHandleManifest() {
	s.Run("when ticker is not in context", func() {
		s.SetupTest()
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
		h := s.handler()
		h.HandleManifest(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Contains(s.w.Body.String(), `"error"`)
	})

	s.Run("when ticker is in context", func() {
		s.SetupTest()
		ticker := storage.Ticker{
			ID:          1,
			Title:       "Test Ticker",
			Description: "A test ticker for unit testing",
			Active:      true,
		}
		s.ctx.Set("ticker", ticker)
		s.ctx.Request = httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
		h := s.handler()
		h.HandleManifest(s.ctx)

		s.Equal(http.StatusOK, s.w.Code)
		s.Equal("application/manifest+json", s.w.Header().Get("Content-Type"))

		var manifest WebAppManifest
		err := json.Unmarshal(s.w.Body.Bytes(), &manifest)
		s.NoError(err)
		s.Equal("Test Ticker", manifest.Name)
		s.Equal("Test Ticker", manifest.ShortName)
		s.Equal("A test ticker for unit testing", manifest.Description)
		s.Equal("/", manifest.StartURL)
		s.Equal("fullscreen", manifest.Display)
		s.Equal("portrait-primary", manifest.Orientation)
		s.Equal("/", manifest.Scope)
		s.Empty(manifest.Icons)
	})
}

func (s *ManifestTestSuite) handler() handler {
	return handler{
		storage: s.store,
		config:  s.cfg,
	}
}

func TestManifestTestSuite(t *testing.T) {
	suite.Run(t, new(ManifestTestSuite))
}
