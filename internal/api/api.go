package api

import (
	"net/http"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/middleware/auth"
	"github.com/systemli/ticker/internal/api/middleware/cors"
	"github.com/systemli/ticker/internal/api/middleware/logger"
	"github.com/systemli/ticker/internal/api/middleware/prometheus"
	"github.com/systemli/ticker/internal/api/middleware/user"
	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

var log = logrus.New().WithField("package", "api")

type handler struct {
	config  config.Config
	storage storage.TickerStorage
	bridges bridge.Bridges
}

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

func API(config config.Config, storage storage.TickerStorage) *gin.Engine {
	handler := handler{
		config:  config,
		storage: storage,
		bridges: bridge.RegisterBridges(config, storage),
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(logger.Logger(config.LogLevel))
	r.Use(gin.Recovery())
	r.Use(cors.NewCORS())
	r.Use(prometheus.NewPrometheus())
	r.Use(limits.RequestSizeLimiter(1024 * 1024 * 10))

	// the jwt middleware
	authMiddleware := auth.AuthMiddleware(storage, config.Secret)
	userMiddleware := user.UserMiddleware(storage)

	admin := r.Group("/v1/admin").Use(authMiddleware.MiddlewareFunc()).Use(userMiddleware)
	{
		admin.GET("/refresh_token", authMiddleware.RefreshHandler)

		admin.GET("/features", handler.GetFeatures)

		admin.GET(`/tickers`, handler.GetTickers)
		admin.GET(`/tickers/:tickerID`, handler.GetTicker)
		admin.POST(`/tickers`, handler.PostTicker)
		admin.PUT(`/tickers/:tickerID`, handler.PutTicker)
		admin.PUT(`/tickers/:tickerID/twitter`, handler.PutTickerTwitter)
		admin.PUT(`/tickers/:tickerID/telegram`, handler.PutTickerTelegram)
		admin.DELETE(`/tickers/:tickerID`, handler.DeleteTicker)
		admin.PUT(`/tickers/:tickerID/reset`, handler.ResetTicker)
		admin.GET(`/tickers/:tickerID/users`, handler.GetTickerUsers)
		admin.PUT(`/tickers/:tickerID/users`, handler.PutTickerUsers)
		admin.DELETE(`/tickers/:tickerID/users/:userID`, handler.DeleteTickerUser)

		admin.GET(`/tickers/:tickerID/messages`, handler.GetMessages)
		admin.GET(`/tickers/:tickerID/messages/:messageID`, handler.GetMessage)
		admin.POST(`/tickers/:tickerID/messages`, handler.PostMessage)
		admin.DELETE(`/tickers/:tickerID/messages/:messageID`, handler.DeleteMessage)

		admin.POST(`/upload`, handler.PostUpload)

		admin.GET(`/users`, handler.GetUsers)
		admin.GET(`/users/:userID`, handler.GetUser)
		admin.POST(`/users`, handler.PostUser)
		admin.PUT(`/users/:userID`, handler.PutUser)
		admin.DELETE(`/users/:userID`, handler.DeleteUser)

		admin.GET(`/settings/:name`, handler.GetSetting)
		admin.PUT(`/settings/inactive_settings`, handler.PutInactiveSettings)
		admin.PUT(`/settings/refresh_interval`, handler.PutRefreshInterval)
	}

	public := r.Group("/v1").Use()
	{
		public.POST(`/admin/login`, authMiddleware.LoginHandler)
		public.POST(`/admin/auth/twitter/request_token`, handler.PostTwitterRequestToken)
		public.POST(`/admin/auth/twitter`, handler.PostAuthTwitter)

		public.GET(`/init`, handler.GetInit)
		public.GET(`/timeline`, handler.GetTimeline)
	}

	r.GET(`/media/:fileName`, handler.GetMedia)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}
