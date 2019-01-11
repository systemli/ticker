package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	reqCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of requests",
	}, []string{"handler", "origin", "code"})

	reqDur = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "http_request_duration_seconds",
		Help: "The HTTP requests latency in seconds",
	}, []string{"handler", "origin", "code"})
)

// NewPrometheus returns the Gin Middleware for collecting basic http metrics.
func NewPrometheus() gin.HandlerFunc {
	err := prometheus.Register(reqCnt)
	if err != nil {
		log.WithError(err).Error(`"reqCnt" could not be registered in Prometheus`)
	}
	err = prometheus.Register(reqDur)
	if err != nil {
		log.WithError(err).Error(`"reqDur" could not be registered in Prometheus`)
	}

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		handler := prepareHandler(c.HandlerName())
		origin := prepareOrigin(c)
		code := strconv.Itoa(c.Writer.Status())

		elapsed := float64(time.Since(start)) / float64(time.Second)
		reqDur.WithLabelValues(handler, origin, code).Observe(elapsed)
		reqCnt.WithLabelValues(handler, origin, code).Inc()
	}
}

func prepareHandler(h string) string {
	s := strings.Split(h, ".")

	return s[len(s)-1]
}

func prepareOrigin(c *gin.Context) string {
	domain, err := GetDomain(c)
	if err != nil {
		return ""
	}

	return domain
}
