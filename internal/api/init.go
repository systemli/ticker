package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

//Settings holds options for frontend settings
type Settings struct {
	RefreshInterval int `json:"refresh_interval,omitempty"`
}

//GetInit returns the basic settings for the ticker.
func GetInit(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")

	u, err := url.Parse(origin)
	if err != nil {
		//TODO: Handle Error
	}

	domain := u.Host
	if strings.HasPrefix(domain, "www.") {
		domain = domain[4:]
	}
	if strings.Contains(domain, ":") {
		parts := strings.Split(domain, ":")
		domain = parts[0]
	}

	var ticker Ticker

	settings := Settings{
		RefreshInterval: 10,
	}

	err = DB.One("Domain", domain, &ticker)
	if err != nil {
		c.JSON(http.StatusOK, JSONResponse{
			Data:   map[string]interface{}{"ticker": nil, "settings": settings},
			Status: ResponseSuccess,
			Error:  nil,
		})
		return
	}

	c.JSON(http.StatusOK, JSONResponse{
		Data:   map[string]interface{}{"ticker": ticker, "settings": settings},
		Status: ResponseSuccess,
		Error:  nil,
	})
	return
}
