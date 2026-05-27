package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetInit(c *gin.Context) {
	settings := response.Settings{}
	origin, err := helper.GetOrigin(c)
	if err != nil {
		c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": nil, "settings": settings}))
		return
	}

	ticker, err := h.stores.Tickers.FindTickerByOrigin(origin)
	if err != nil || !ticker.Active {
		settings.InactiveSettings = storage.GetSettings(h.stores.Settings, storage.InactiveSetting)
		c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": nil, "settings": settings}))
		return
	}

	data := map[string]interface{}{"ticker": response.InitTickerResponse(ticker), "settings": settings}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}
