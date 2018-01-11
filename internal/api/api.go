package api

import "github.com/gin-gonic/gin"

//Returns the Gin Engine
func API() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1").Use()
	{
		// Endpoints for tickers
		//TODO: Authentication
		v1.GET(`/admin/tickers`, GetTickers)
		v1.GET(`/admin/tickers/:id`, GetTicker)
		v1.POST(`/admin/tickers`, PostTicker)
		v1.PUT(`/admin/tickers/:id`, PutTicker)
		v1.DELETE(`/admin/tickers/:id`, DeleteTicker)

		// Endpoints for messages
		//TODO: Authentication
		v1.GET(`/admin/messages`, GetMessages)
		v1.GET(`/admin/messages/:id`, GetMessage)
		v1.POST(`/admin/messages`, PostMessage)
		v1.PUT(`/admin/messages/:id`, PutMessage)
		v1.DELETE(`/admin/messages/:id`, DeleteMessage)
	}

	return r
}
