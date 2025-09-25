package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		fields := logrus.Fields{
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"status":         c.Writer.Status(),
			"latency":        latency.String(),
			"correlation_id": c.GetString("correlation_id"),
		}
		if userID := c.GetString("user_id"); userID != "" {
			fields["user_id"] = userID
		}

		// Prometheus metrics
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(latency.Seconds())
		logrus.WithFields(fields).Info("request completed")
	}
}
