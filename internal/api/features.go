package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type FeaturesResponse map[string]bool

func NewFeaturesResponse(config config.Config, storage storage.Storage) FeaturesResponse {
	telegramSettings := storage.GetTelegramSettings()
	return FeaturesResponse{
		"telegramEnabled":    telegramSettings.Token != "",
		"signalGroupEnabled": config.SignalGroup.Enabled(),
	}
}

func (h *handler) GetFeatures(c *gin.Context) {
	features := NewFeaturesResponse(h.config, h.storage)
	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"features": features}))
}
