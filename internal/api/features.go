package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

type FeaturesResponse map[string]bool

func NewFeaturesResponse(settings storage.SettingsStore) FeaturesResponse {
	telegramSettings := storage.GetSettings(settings, storage.TelegramSetting)
	signalGroupSettings := storage.GetSettings(settings, storage.SignalGroupSetting)
	return FeaturesResponse{
		"telegramEnabled":    telegramSettings.Token != "",
		"signalGroupEnabled": signalGroupSettings.Enabled(),
	}
}

func (h *handler) GetFeatures(c *gin.Context) {
	features := NewFeaturesResponse(h.stores.Settings)
	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"features": features}))
}
