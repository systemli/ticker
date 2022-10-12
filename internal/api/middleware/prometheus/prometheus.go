package prometheus

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/helper"
)

var (
	requestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "The HTTP requests latency in seconds",
		Buckets: []float64{0.5, 0.9, 0.95, 0.99},
	}, []string{"handler", "origin", "code"})
)

func NewPrometheus() gin.HandlerFunc {
	err := prometheus.Register(requestDurationHistogram)
	if err != nil {
		log.WithError(err).Error(`"requestDurationHistogram" could not be registered in Prometheus`)
		return func(c *gin.Context) {}
	}

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		handler := prepareHandler(c.HandlerName())
		origin := prepareOrigin(c)
		code := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)

		requestDurationHistogram.WithLabelValues(handler, origin, code).Observe(elapsed)
	}
}

func prepareHandler(h string) string {
	s := strings.Split(h, ".")

	return strings.TrimSuffix(s[len(s)-1], "-fm")
}

func prepareOrigin(c *gin.Context) string {
	domain, err := helper.GetDomain(c)
	if err != nil {
		return ""
	}

	return domain
}
