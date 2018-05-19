package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"strings"
	"net/url"
)

//Returns the Gin Engine
func API() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{`GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`}

	r.Use(cors.New(config))

	v1 := r.Group("/v1").Use()
	{
		// Endpoints for tickers
		//TODO: Authentication
		v1.GET(`/admin/tickers`, GetTickers)
		v1.GET(`/admin/tickers/:tickerID`, GetTicker)
		v1.POST(`/admin/tickers`, PostTicker)
		v1.PUT(`/admin/tickers/:tickerID`, PutTicker)
		v1.DELETE(`/admin/tickers/:tickerID`, DeleteTicker)

		v1.GET(`/admin/tickers/:tickerID/messages`, GetMessages)
		v1.GET(`/admin/tickers/:tickerID/messages/:messageID`, GetMessage)
		v1.POST(`/admin/tickers/:tickerID/messages`, PostMessage)
		v1.DELETE(`/admin/tickers/:tickerID/messages/:messageID`, DeleteMessage)

		v1.GET(`/init`, GetInit)
		v1.GET(`/timeline`, GetTimeline)
	}

	return r
}

//
func GetDomain(c *gin.Context) (string, error) {
	origin := c.Request.Header.Get("Origin")

	u, err := url.Parse(origin)
	if err != nil {
		return "", err
	}

	domain := u.Host
	if strings.HasPrefix(domain, "www.") {
		domain = domain[4:]
	}
	if strings.Contains(domain, ":") {
		parts := strings.Split(domain, ":")
		domain = parts[0]
	}

	return domain, nil
}