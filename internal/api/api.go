package api

import (
	"net/http"
	"time"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/middleware/auth"
	"github.com/systemli/ticker/internal/api/middleware/cors"
	"github.com/systemli/ticker/internal/api/middleware/logger"
	"github.com/systemli/ticker/internal/api/middleware/me"
	"github.com/systemli/ticker/internal/api/middleware/message"
	"github.com/systemli/ticker/internal/api/middleware/prometheus"
	"github.com/systemli/ticker/internal/api/middleware/response_cache"
	"github.com/systemli/ticker/internal/api/middleware/ticker"
	"github.com/systemli/ticker/internal/api/middleware/user"
	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/cache"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

var log = logrus.New().WithField("package", "api")

type handler struct {
	config  config.Config
	storage storage.Storage
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

func API(config config.Config, store storage.Storage, log *logrus.Logger) *gin.Engine {
	handler := handler{
		config:  config,
		storage: store,
		bridges: bridge.RegisterBridges(config, store),
	}

	// TOOD: Make this configurable via config file
	cacheTtl := 30 * time.Second
	inMemoryCache := cache.NewCache(5 * time.Minute)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(logger.Logger(log))
	r.Use(gin.Recovery())
	r.Use(cors.NewCORS())
	r.Use(prometheus.NewPrometheus())
	r.Use(limits.RequestSizeLimiter(1024 * 1024 * 10))

	// the jwt middleware
	authMiddleware := auth.AuthMiddleware(store, config.Secret)

	admin := r.Group("/v1/admin")
	{
		meMiddleware := me.MeMiddleware(store)
		admin.Use(authMiddleware.MiddlewareFunc())
		admin.Use(meMiddleware)

		admin.GET("/refresh_token", authMiddleware.RefreshHandler)

		admin.GET("/features", handler.GetFeatures)

		admin.GET(`/tickers`, handler.GetTickers)
		admin.GET(`/tickers/:tickerID`, ticker.PrefetchTicker(store, storage.WithPreload()), handler.GetTicker)
		admin.POST(`/tickers`, user.NeedAdmin(), handler.PostTicker)
		admin.PUT(`/tickers/:tickerID`, ticker.PrefetchTicker(store, storage.WithPreload()), handler.PutTicker)
		admin.PUT(`/tickers/:tickerID/telegram`, ticker.PrefetchTicker(store, storage.WithPreload()), handler.PutTickerTelegram)
		admin.DELETE(`/tickers/:tickerID/telegram`, ticker.PrefetchTicker(store, storage.WithPreload()), handler.DeleteTickerTelegram)
		admin.PUT(`/tickers/:tickerID/mastodon`, ticker.PrefetchTicker(store, storage.WithPreload()), handler.PutTickerMastodon)
		admin.DELETE(`/tickers/:tickerID/mastodon`, ticker.PrefetchTicker(store, storage.WithPreload()), handler.DeleteTickerMastodon)
		admin.DELETE(`/tickers/:tickerID`, user.NeedAdmin(), ticker.PrefetchTicker(store), handler.DeleteTicker)
		admin.PUT(`/tickers/:tickerID/reset`, ticker.PrefetchTicker(store, storage.WithPreload()), ticker.PrefetchTicker(store), handler.ResetTicker)
		admin.GET(`/tickers/:tickerID/users`, ticker.PrefetchTicker(store), handler.GetTickerUsers)
		admin.PUT(`/tickers/:tickerID/users`, user.NeedAdmin(), ticker.PrefetchTicker(store), handler.PutTickerUsers)
		admin.DELETE(`/tickers/:tickerID/users/:userID`, user.NeedAdmin(), ticker.PrefetchTicker(store), handler.DeleteTickerUser)

		admin.GET(`/tickers/:tickerID/messages`, ticker.PrefetchTicker(store, storage.WithPreload()), handler.GetMessages)
		admin.GET(`/tickers/:tickerID/messages/:messageID`, ticker.PrefetchTicker(store, storage.WithPreload()), message.PrefetchMessage(store), handler.GetMessage)
		admin.POST(`/tickers/:tickerID/messages`, ticker.PrefetchTicker(store), handler.PostMessage)
		admin.DELETE(`/tickers/:tickerID/messages/:messageID`, ticker.PrefetchTicker(store), message.PrefetchMessage(store), handler.DeleteMessage)

		admin.POST(`/upload`, handler.PostUpload)

		admin.GET(`/users`, user.NeedAdmin(), handler.GetUsers)
		admin.GET(`/users/:userID`, user.PrefetchUser(store), handler.GetUser)
		admin.POST(`/users`, user.NeedAdmin(), handler.PostUser)
		admin.PUT(`/users/me`, handler.PutMe)
		admin.PUT(`/users/:userID`, user.NeedAdmin(), user.PrefetchUser(store), handler.PutUser)
		admin.DELETE(`/users/:userID`, user.NeedAdmin(), user.PrefetchUser(store), handler.DeleteUser)

		admin.GET(`/settings/:name`, user.NeedAdmin(), handler.GetSetting)
		admin.PUT(`/settings/inactive_settings`, user.NeedAdmin(), handler.PutInactiveSettings)
		admin.PUT(`/settings/refresh_interval`, user.NeedAdmin(), handler.PutRefreshInterval)
	}

	public := r.Group("/v1").Use()
	{
		public.POST(`/admin/login`, authMiddleware.LoginHandler)

		public.GET(`/init`, response_cache.CachePage(inMemoryCache, cacheTtl, handler.GetInit))
		public.GET(`/timeline`, ticker.PrefetchTickerFromRequest(store), response_cache.CachePage(inMemoryCache, cacheTtl, handler.GetTimeline))
		public.GET(`/feed`, ticker.PrefetchTickerFromRequest(store), response_cache.CachePage(inMemoryCache, cacheTtl, handler.GetFeed))
	}

	r.GET(`/media/:fileName`, handler.GetMedia)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}
