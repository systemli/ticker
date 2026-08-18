package api

import (
	"net/http"
	"time"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/middleware/auth"
	"github.com/systemli/ticker/internal/api/middleware/cors"
	loggerMiddleware "github.com/systemli/ticker/internal/api/middleware/logger"
	"github.com/systemli/ticker/internal/api/middleware/me"
	"github.com/systemli/ticker/internal/api/middleware/message"
	"github.com/systemli/ticker/internal/api/middleware/prometheus"
	"github.com/systemli/ticker/internal/api/middleware/response_cache"
	"github.com/systemli/ticker/internal/api/middleware/ticker"
	"github.com/systemli/ticker/internal/api/middleware/user"
	"github.com/systemli/ticker/internal/api/realtime"
	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/cache"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/logger"
	"github.com/systemli/ticker/internal/storage"
)

var log = logger.GetWithPackage("api")

// Server wraps the gin engine and realtime engine for graceful shutdown
type Server struct {
	Router   *gin.Engine
	Realtime *realtime.Engine
}

type handler struct {
	config   config.Config
	stores   storage.Stores
	bridges  bridge.Bridges
	cache    *cache.Cache
	realtime *realtime.Engine
}

func API(config config.Config, stores storage.Stores) *Server {
	inMemoryCache := cache.NewCache(5 * time.Minute)

	ws := realtime.New()
	go ws.Run()

	handler := handler{
		config:   config,
		stores:   stores,
		bridges:  bridge.RegisterBridges(config, stores),
		cache:    inMemoryCache,
		realtime: ws,
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(loggerMiddleware.Logger(log.Logger))
	r.Use(gin.Recovery())
	r.Use(cors.NewCORS())
	r.Use(prometheus.NewPrometheus())
	r.Use(limits.RequestSizeLimiter(1024 * 1024 * 10))

	// the jwt middleware
	authMiddleware := auth.AuthMiddleware(stores.Users, config.Secret)

	admin := r.Group("/v1/admin")
	{
		meMiddleware := me.MeMiddleware(stores.Users)
		admin.Use(authMiddleware.MiddlewareFunc())
		admin.Use(meMiddleware)

		admin.GET("/refresh_token", authMiddleware.RefreshHandler)

		admin.GET("/features", handler.GetFeatures)

		admin.GET(`/tickers`, handler.GetTickers)
		admin.GET(`/tickers/:tickerID`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.GetTicker)
		admin.POST(`/tickers`, user.NeedAdmin(), handler.PostTicker)
		admin.PUT(`/tickers/:tickerID`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.PutTicker)
		admin.DELETE(`/tickers/:tickerID/websites`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.DeleteTickerWebsites)
		admin.PUT(`/tickers/:tickerID/websites`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.PutTickerWebsites)
		admin.PUT(`/tickers/:tickerID/telegram`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.PutTickerTelegram)
		admin.DELETE(`/tickers/:tickerID/telegram`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.DeleteTickerTelegram)
		admin.PUT(`/tickers/:tickerID/mastodon`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.PutTickerMastodon)
		admin.DELETE(`/tickers/:tickerID/mastodon`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.DeleteTickerMastodon)
		admin.PUT(`/tickers/:tickerID/bluesky`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.PutTickerBluesky)
		admin.DELETE(`/tickers/:tickerID/bluesky`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.DeleteTickerBluesky)
		admin.PUT(`/tickers/:tickerID/signal_group`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.PutTickerSignalGroup)
		admin.DELETE(`/tickers/:tickerID/signal_group`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.DeleteTickerSignalGroup)
		admin.PUT(`/tickers/:tickerID/signal_group/admin`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.PutTickerSignalGroupAdmin)
		admin.DELETE(`/tickers/:tickerID`, user.NeedAdmin(), ticker.PrefetchTicker(stores.Tickers), handler.DeleteTicker)
		admin.PUT(`/tickers/:tickerID/reset`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), ticker.PrefetchTicker(stores.Tickers), handler.ResetTicker)
		admin.GET(`/tickers/:tickerID/users`, ticker.PrefetchTicker(stores.Tickers), handler.GetTickerUsers)
		admin.PUT(`/tickers/:tickerID/users`, user.NeedAdmin(), ticker.PrefetchTicker(stores.Tickers), handler.PutTickerUsers)
		admin.DELETE(`/tickers/:tickerID/users/:userID`, user.NeedAdmin(), ticker.PrefetchTicker(stores.Tickers), handler.DeleteTickerUser)

		admin.GET(`/tickers/:tickerID/messages`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.GetMessages)
		admin.GET(`/tickers/:tickerID/messages/:messageID`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), message.PrefetchMessage(stores.Messages), handler.GetMessage)
		admin.POST(`/tickers/:tickerID/messages`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), handler.PostMessage)
		admin.DELETE(`/tickers/:tickerID/messages/:messageID`, ticker.PrefetchTicker(stores.Tickers, storage.WithPreload()), message.PrefetchMessage(stores.Messages), handler.DeleteMessage)

		admin.POST(`/upload`, handler.PostUpload)

		admin.GET(`/users`, user.NeedAdmin(), handler.GetUsers)
		admin.GET(`/users/:userID`, user.PrefetchUser(stores.Users), handler.GetUser)
		admin.POST(`/users`, user.NeedAdmin(), handler.PostUser)
		admin.PUT(`/users/me`, handler.PutMe)
		admin.PUT(`/users/:userID`, user.NeedAdmin(), user.PrefetchUser(stores.Users), handler.PutUser)
		admin.DELETE(`/users/:userID`, user.NeedAdmin(), user.PrefetchUser(stores.Users), handler.DeleteUser)

		admin.GET(`/settings/:name`, user.NeedAdmin(), handler.GetSetting)
		admin.PUT(`/settings/inactive_settings`, user.NeedAdmin(), handler.PutInactiveSettings)
		admin.PUT(`/settings/telegram_settings`, user.NeedAdmin(), handler.PutTelegramSettings)
		admin.PUT(`/settings/signal_group_settings`, user.NeedAdmin(), handler.PutSignalGroupSettings)
	}

	public := r.Group("/v1").Use()
	{
		public.POST(`/admin/login`, authMiddleware.LoginHandler)

		public.GET(`/init`, response_cache.CachePage(inMemoryCache, 5*time.Minute, handler.GetInit))
		public.GET(`/manifest.json`, ticker.PrefetchTickerFromRequest(stores.Tickers), handler.HandleManifest)
		public.GET(`/timeline`, ticker.PrefetchTickerFromRequest(stores.Tickers), response_cache.CachePage(inMemoryCache, 10*time.Second, handler.GetTimeline))
		public.GET(`/feed`, ticker.PrefetchTickerFromRequest(stores.Tickers), response_cache.CachePage(inMemoryCache, 5*time.Minute, handler.GetFeed))
		public.GET(`/ws`, ticker.PrefetchTickerFromRequest(stores.Tickers), handler.HandleWebSocket)
	}

	r.GET(`/media/:fileName`, handler.GetMedia)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return &Server{
		Router:   r,
		Realtime: ws,
	}
}
