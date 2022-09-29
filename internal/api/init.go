package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
)

// GetInit returns the basic settings for the ticker.
// @Summary      Retrieves the initial ticker configuration
// @Description  The first request for retrieving information about the ticker. It is mandatory that the browser sends
// @Description  the origin as a header. This can be overwritten with a query parameter.
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        origin  query     string  false  "Origin from the ticker, e.g. demoticker.org"
// @Param        origin  header    string  false  "Origin from the ticker, e.g. http://demoticker.org"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /init [get]
func (h *handler) GetInit(c *gin.Context) {
	settings := response.Settings{
		RefreshInterval: h.storage.GetRefreshIntervalSetting().Value.(int),
	}
	domain, err := helper.GetDomain(c)
	if err != nil {
		c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": nil, "settings": settings}))
		return
	}

	ticker, err := h.storage.FindTickerByDomain(domain)
	if err != nil || !ticker.Active {
		settings.InactiveSettings = h.storage.GetInactiveSetting().Value
		c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"ticker": nil, "settings": settings}))
		return
	}

	data := map[string]interface{}{"ticker": response.TickerResponse(ticker, h.config), "settings": settings}
	c.JSON(http.StatusOK, response.SuccessResponse(data))
}
