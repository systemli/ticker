package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
)

type FeaturesResponse map[string]bool

func NewFeaturesResponse(config config.Config) FeaturesResponse {
	return FeaturesResponse{
		"telegramEnabled": config.TelegramEnabled(),
	}
}

func (h *handler) GetFeatures(c *gin.Context) {
	features := NewFeaturesResponse(h.config)
	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"features": features}))
}
