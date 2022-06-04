package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
)

type FeaturesResponse map[string]bool

func NewFeaturesResponse() *FeaturesResponse {
	return &FeaturesResponse{
		"twitter_enabled":  Config.TwitterEnabled(),
		"telegram_enabled": Config.TelegramEnabled(),
	}
}

func GetFeaturesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, NewJSONSuccessResponse("features", NewFeaturesResponse()))
}
