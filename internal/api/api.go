package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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
	}

	return r
}
