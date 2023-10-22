package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
)

func (h *handler) GetInit(c *gin.Context) {
	settings := response.Settings{
		RefreshInterval: h.storage.GetRefreshIntervalSettings().RefreshInterval,
	}
	domain, err := helper.GetDomain(c)
	if err != nil {
		c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": nil, "settings": settings}))
		return
	}

	ticker, err := h.storage.FindTickerByDomain(domain)
	if err != nil || !ticker.Active {
		settings.InactiveSettings = h.storage.GetInactiveSettings()
		c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": nil, "settings": settings}))
		return
	}

	data := map[string]interface{}{"ticker": response.InitTickerResponse(ticker), "settings": settings}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}
