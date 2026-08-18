package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type FeaturesTestSuite struct {
	suite.Suite
}

func (s *FeaturesTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *FeaturesTestSuite) TestGetFeatures() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	store := storage.NewMockStorage()

	// Mock GetTelegramSettings to return empty token (disabled)
	store.Settings.MockGetTelegram(storage.TelegramSettings{Token: ""})
	// Mock GetSignalGroupSettings to return empty settings (disabled)
	store.Settings.MockGetSignalGroup(storage.DefaultSignalGroupSettings())

	h := handler{
		stores: store.Stores(),
		config: config.LoadConfig(""),
	}

	h.GetFeatures(c)

	s.Equal(http.StatusOK, w.Code)
	s.Equal(`{"data":{"features":{"signalGroupEnabled":false,"telegramEnabled":false}},"status":"success","error":{}}`, w.Body.String())
}

func TestFeaturesTestSuite(t *testing.T) {
	suite.Run(t, new(FeaturesTestSuite))
}
