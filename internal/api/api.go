package api

import (
	"net/url"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//Returns the Gin Engine
func API() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{`*`}
	config.AllowMethods = []string{`GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`}

	r.Use(cors.New(config))

	// the jwt middleware
	authMiddleware := AuthMiddleware()

	admin := r.Group("/v1/admin").Use(authMiddleware.MiddlewareFunc())
	{
		admin.GET("/refresh_token", authMiddleware.RefreshHandler)

		admin.GET(`/tickers`, GetTickers)
		admin.GET(`/tickers/:tickerID`, GetTicker)
		admin.POST(`/tickers`, PostTicker)
		admin.PUT(`/tickers/:tickerID`, PutTicker)
		admin.DELETE(`/tickers/:tickerID`, DeleteTicker)

		admin.GET(`/tickers/:tickerID/messages`, GetMessages)
		admin.GET(`/tickers/:tickerID/messages/:messageID`, GetMessage)
		admin.POST(`/tickers/:tickerID/messages`, PostMessage)
		admin.DELETE(`/tickers/:tickerID/messages/:messageID`, DeleteMessage)

		admin.GET(`/users`, GetUsers)
		admin.GET(`/users/:userID`, GetUser)
		admin.POST(`/users`, PostUser)
		admin.PUT(`/users/:userID`, PutUser)
		admin.DELETE(`/users/:userID`, DeleteUser)
	}

	public := r.Group("/v1").Use()
	{
		public.POST(`/admin/login`, authMiddleware.LoginHandler)
		public.GET(`/init`, GetInit)
		public.GET(`/timeline`, GetTimeline)
	}

	return r
}

//
func GetDomain(c *gin.Context) (string, error) {
	origin := c.Request.Header.Get("Origin")

	if origin == "" {
		return "", errors.New("Origin header not found")
	}

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
