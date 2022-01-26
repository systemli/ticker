package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

//Settings represents the general Settings for TickerResponse in Init Response.
type Settings struct {
	RefreshInterval  int         `json:"refresh_interval,omitempty"`
	InactiveSettings interface{} `json:"inactive_settings,omitempty"`
}

type initResponse struct {
	Data   initData    `json:"data"`
	Status string      `json:"status"`
	Error  interface{} `json:"error"`
}

type initData struct {
	Ticker   *TickerResponse `json:"ticker"`
	Settings *Settings       `json:"settings"`
}

// GetInitHandler returns the basic settings for the ticker.
// @Summary      Retrieves the initial ticker configuration
// @Description  The first request for retrieving information about the ticker. It is mandatory that the browser sends
// @Description  the origin as a header. This can be overwritten with a query parameter.
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        origin  query     string  false  "Origin from the ticker, e.g. demoticker.org"
// @Param        origin  header    string  false  "Origin from the ticker, e.g. http://demoticker.org"
// @Success      200     {object}  initResponse
// @Failure      500     {object}  interface{}
// @Router       /init [get]
func GetInitHandler(c *gin.Context) {
	settings := settings()
	domain, err := GetDomain(c)
	if err != nil {
		emptyTickerResponse(c, settings)
		return
	}

	ticker, err := FindTicker(domain)
	if err != nil || !ticker.Active {
		settings.InactiveSettings = GetInactiveSettings().Value
		emptyTickerResponse(c, settings)
		return
	}

	c.JSON(http.StatusOK, initResponse{
		Data:   initData{Ticker: NewTickerResponse(ticker), Settings: settings},
		Status: ResponseSuccess,
		Error:  nil,
	})
}

func settings() *Settings {
	return &Settings{
		RefreshInterval: GetRefreshIntervalValue(),
	}
}

func emptyTickerResponse(c *gin.Context, s *Settings) {
	c.JSON(http.StatusOK, initResponse{
		Data:   initData{Ticker: nil, Settings: s},
		Status: ResponseSuccess,
		Error:  nil,
	})
}
