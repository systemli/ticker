package prometheus

import (
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/helper"
)

var (
	requestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "The HTTP requests latency in seconds",
	}, []string{"method", "path", "origin", "code"})
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

		method := c.Request.Method
		path := c.FullPath()
		origin := prepareOrigin(c)
		code := strconv.Itoa(c.Writer.Status())

		requestDurationHistogram.WithLabelValues(method, path, origin, code).Observe(time.Since(start).Seconds())
	}
}

func prepareOrigin(c *gin.Context) string {
	origin, err := helper.GetOrigin(c)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`^https?://`)

	return re.ReplaceAllString(origin, "")
}
