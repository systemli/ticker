package api

import (
	"net/http"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/middleware/auth"
	"github.com/systemli/ticker/internal/api/middleware/cors"
	"github.com/systemli/ticker/internal/api/middleware/logger"
	"github.com/systemli/ticker/internal/api/middleware/me"
	"github.com/systemli/ticker/internal/api/middleware/message"
	"github.com/systemli/ticker/internal/api/middleware/prometheus"
	"github.com/systemli/ticker/internal/api/middleware/ticker"
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

	admin := r.Group("/v1/admin")
	{
		meMiddleware := me.MeMiddleware(storage)
		admin.Use(authMiddleware.MiddlewareFunc())
		admin.Use(meMiddleware)

		admin.GET("/refresh_token", authMiddleware.RefreshHandler)

		admin.GET("/features", handler.GetFeatures)

		admin.GET(`/tickers`, handler.GetTickers)
		admin.GET(`/tickers/:tickerID`, ticker.PrefetchTicker(storage), handler.GetTicker)
		admin.POST(`/tickers`, user.NeedAdmin(), handler.PostTicker)
		admin.PUT(`/tickers/:tickerID`, ticker.PrefetchTicker(storage), handler.PutTicker)
		admin.PUT(`/tickers/:tickerID/twitter`, ticker.PrefetchTicker(storage), handler.PutTickerTwitter)
		admin.PUT(`/tickers/:tickerID/telegram`, ticker.PrefetchTicker(storage), handler.PutTickerTelegram)
		admin.DELETE(`/tickers/:tickerID`, user.NeedAdmin(), ticker.PrefetchTicker(storage), handler.DeleteTicker)
		admin.PUT(`/tickers/:tickerID/reset`, ticker.PrefetchTicker(storage), ticker.PrefetchTicker(storage), handler.ResetTicker)
		admin.GET(`/tickers/:tickerID/users`, ticker.PrefetchTicker(storage), handler.GetTickerUsers)
		admin.PUT(`/tickers/:tickerID/users`, user.NeedAdmin(), ticker.PrefetchTicker(storage), handler.PutTickerUsers)
		admin.DELETE(`/tickers/:tickerID/users/:userID`, user.NeedAdmin(), ticker.PrefetchTicker(storage), handler.DeleteTickerUser)

		admin.GET(`/tickers/:tickerID/messages`, ticker.PrefetchTicker(storage), handler.GetMessages)
		admin.GET(`/tickers/:tickerID/messages/:messageID`, ticker.PrefetchTicker(storage), message.PrefetchMessage(storage), handler.GetMessage)
		admin.POST(`/tickers/:tickerID/messages`, ticker.PrefetchTicker(storage), handler.PostMessage)
		admin.DELETE(`/tickers/:tickerID/messages/:messageID`, ticker.PrefetchTicker(storage), message.PrefetchMessage(storage), handler.DeleteMessage)

		admin.POST(`/upload`, handler.PostUpload)

		admin.GET(`/users`, user.NeedAdmin(), handler.GetUsers)
		admin.GET(`/users/:userID`, user.PrefetchUser(storage), handler.GetUser)
		admin.POST(`/users`, user.NeedAdmin(), handler.PostUser)
		admin.PUT(`/users/:userID`, user.NeedAdmin(), user.PrefetchUser(storage), handler.PutUser)
		admin.DELETE(`/users/:userID`, user.NeedAdmin(), user.PrefetchUser(storage), handler.DeleteUser)

		admin.GET(`/settings/:name`, user.NeedAdmin(), handler.GetSetting)
		admin.PUT(`/settings/inactive_settings`, user.NeedAdmin(), handler.PutInactiveSettings)
		admin.PUT(`/settings/refresh_interval`, user.NeedAdmin(), handler.PutRefreshInterval)
	}

	public := r.Group("/v1").Use()
	{
		public.POST(`/admin/login`, authMiddleware.LoginHandler)
		public.POST(`/admin/auth/twitter/request_token`, handler.PostTwitterRequestToken)
		public.POST(`/admin/auth/twitter`, handler.PostAuthTwitter)

		public.GET(`/init`, handler.GetInit)
		public.GET(`/timeline`, ticker.PrefetchTickerFromRequest(storage), handler.GetTimeline)
	}

	r.GET(`/media/:fileName`, handler.GetMedia)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}
