package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := c.GetHeader("X-Correlation-ID")
		if cid == "" {
			cid = uuid.New().String()
		}
		c.Set("correlation_id", cid)
		c.Writer.Header().Set("X-Correlation-ID", cid)
		c.Next()
	}
}
