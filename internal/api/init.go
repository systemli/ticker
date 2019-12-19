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

//GetInitHandler returns the basic settings for the ticker.
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

	c.JSON(http.StatusOK, JSONResponse{
		//TODO: Build NewTickerPublicResponse to hide unnecessary information
		Data:   map[string]interface{}{"ticker": NewTickerResponse(ticker), "settings": settings},
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
	c.JSON(http.StatusOK, JSONResponse{
		Data:   map[string]interface{}{"ticker": nil, "settings": s},
		Status: ResponseSuccess,
		Error:  nil,
	})
}
