package api

import (
	"net/url"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/toorop/gin-logrus"
	"github.com/sirupsen/logrus"

	"git.codecoop.org/systemli/ticker/internal/model"
)

//Returns the Gin Engine
func API() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(Logger())
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Origin", "Content-Length", "Content-Type"}
	config.AllowMethods = []string{`GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`}

	r.Use(cors.New(config))

	// the jwt middleware
	authMiddleware := AuthMiddleware()

	admin := r.Group("/v1/admin").Use(authMiddleware.MiddlewareFunc()).Use(UserMiddleware())
	{
		admin.GET("/refresh_token", authMiddleware.RefreshHandler)

		admin.GET(`/tickers`, GetTickersHandler)
		admin.GET(`/tickers/:tickerID`, GetTickerHandler)
		admin.POST(`/tickers`, PostTickerHandler)
		admin.PUT(`/tickers/:tickerID`, PutTickerHandler)
		admin.PUT(`/tickers/:tickerID/twitter`, PutTickerTwitterHandler)
		admin.DELETE(`/tickers/:tickerID`, DeleteTickerHandler)

		admin.GET(`/tickers/:tickerID/messages`, GetMessagesHandler)
		admin.GET(`/tickers/:tickerID/messages/:messageID`, GetMessageHandler)
		admin.POST(`/tickers/:tickerID/messages`, PostMessageHandler)
		admin.DELETE(`/tickers/:tickerID/messages/:messageID`, DeleteMessageHandler)

		admin.GET(`/users`, GetUsersHandler)
		admin.GET(`/users/:userID`, GetUserHandler)
		admin.POST(`/users`, PostUserHandler)
		admin.PUT(`/users/:userID`, PutUserHandler)
		admin.DELETE(`/users/:userID`, DeleteUserHandler)

		admin.GET(`/settings/:name`, GetSettingHandler)
		admin.PUT(`/settings/inactive_settings`, PutInactiveSettingsHandler)
	}

	public := r.Group("/v1").Use()
	{
		public.POST(`/admin/login`, authMiddleware.LoginHandler)
		public.POST(`/admin/auth/twitter/request_token`, PostTwitterRequestTokenHandler)
		public.POST(`/admin/auth/twitter`, PostAuthTwitterHandler)

		public.GET(`/init`, GetInitHandler)
		public.GET(`/timeline`, GetTimelineHandler)
	}

	return r
}

func Logger() gin.HandlerFunc {
	lvl, _ := logrus.ParseLevel(model.Config.LogLevel)
	logger := logrus.New()
	logger.SetLevel(lvl)

	return ginlogrus.Logger(logger)
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

func Me(c *gin.Context) (model.User, error) {
	var user model.User
	u, exists := c.Get(UserKey)
	if !exists {
		return user, errors.New(model.ErrorUserNotFound)
	}

	return u.(model.User), nil
}

func IsAdmin(c *gin.Context) bool {
	u, err := Me(c)
	if err != nil {
		return false
	}

	return u.IsSuperAdmin
}
