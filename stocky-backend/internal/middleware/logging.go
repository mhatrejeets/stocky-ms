package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
		logrus.WithFields(fields).Info("request completed")
	}
}
