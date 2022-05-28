package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"

	"github.com/systemli/ticker/internal/model"
)

// @title         Ticker API
// @version       1.0
// @description   Service to distribute short messages in support of events, demonstrations, or other time-sensitive events.

// @contact.name  Systemli Admin Team
// @contact.url   https://www.systemli.org/en/contact/
// @contact.email admin@systemli.org

// @license.name  GPLv3
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html

// @host          localhost:8080
// @BasePath      /v1

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

	r.Use(NewPrometheus())

	r.Use(limits.RequestSizeLimiter(1024 * 1024 * 10))

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
		admin.PUT(`/tickers/:tickerID/telegram`, PutTickerTelegramHandler)
		admin.DELETE(`/tickers/:tickerID`, DeleteTickerHandler)
		admin.PUT(`/tickers/:tickerID/reset`, ResetTickerHandler)
		admin.GET(`/tickers/:tickerID/users`, GetTickerUsersHandler)
		admin.PUT(`/tickers/:tickerID/users`, PutTickerUsersHandler)
		admin.DELETE(`/tickers/:tickerID/users/:userID`, DeleteTickerUserHandler)

		admin.GET(`/tickers/:tickerID/messages`, GetMessagesHandler)
		admin.GET(`/tickers/:tickerID/messages/:messageID`, GetMessageHandler)
		admin.POST(`/tickers/:tickerID/messages`, PostMessageHandler)
		admin.DELETE(`/tickers/:tickerID/messages/:messageID`, DeleteMessageHandler)

		admin.POST(`/upload`, PostUpload)

		admin.GET(`/users`, GetUsersHandler)
		admin.GET(`/users/:userID`, GetUserHandler)
		admin.POST(`/users`, PostUserHandler)
		admin.PUT(`/users/:userID`, PutUserHandler)
		admin.DELETE(`/users/:userID`, DeleteUserHandler)

		admin.GET(`/settings/:name`, GetSettingHandler)
		admin.PUT(`/settings/inactive_settings`, PutInactiveSettingsHandler)
		admin.PUT(`/settings/refresh_interval`, PutRefreshIntervalHandler)
	}

	public := r.Group("/v1").Use()
	{
		public.POST(`/admin/login`, authMiddleware.LoginHandler)
		public.POST(`/admin/auth/twitter/request_token`, PostTwitterRequestTokenHandler)
		public.POST(`/admin/auth/twitter`, PostAuthTwitterHandler)

		public.GET(`/init`, GetInitHandler)
		public.GET(`/timeline`, GetTimelineHandler)
	}

	r.GET(`/media/:fileName`, GetMedia)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}

func Logger() gin.HandlerFunc {
	lvl, _ := logrus.ParseLevel(model.Config.LogLevel)
	logger := logrus.New()
	logger.SetLevel(lvl)

	return ginlogrus.Logger(logger)
}
