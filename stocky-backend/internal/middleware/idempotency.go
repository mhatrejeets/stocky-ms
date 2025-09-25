package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhatrejeets/stocky-ms/internal/repo"
)

func Idempotency(repo repo.RewardRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost && c.FullPath() == "/api/v1/reward" {
			key := c.GetHeader("Idempotency-Key")
			if key != "" {
				exists, resp := repo.CheckIdempotencyKey(c.Request.Context(), key)
				if exists {
					c.JSON(http.StatusOK, resp)
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}
