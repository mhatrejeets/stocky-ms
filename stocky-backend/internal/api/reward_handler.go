package api

import (
	"context"
	"net/http"

	"github.com/mhatrejeets/stocky-ms/internal/model"
	"github.com/mhatrejeets/stocky-ms/internal/service"

	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	Service *service.RewardService
}

func (h *RewardHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/reward", h.CreateReward)
		v1.GET("/today-stocks/:userId", h.GetTodayStocks)
		v1.GET("/historical-inr/:userId", h.GetHistoricalINR)
		v1.GET("/stats/:userId", h.GetStats)
		v1.GET("/portfolio/:userId", h.GetPortfolio)
	}
}

func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req model.CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
		return
	}
	idempotencyKey := c.GetHeader("Idempotency-Key")
	result := h.Service.CreateReward(context.Background(), userID, req, idempotencyKey)
	if result.Conflict {
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate reward"})
		return
	}
	if result.Err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "reward_id": result.RewardID})
}

func (h *RewardHandler) GetTodayStocks(c *gin.Context) {
	userID := c.Param("userId")
	// TODO: Call service/repo to get today's rewards
	c.JSON(http.StatusOK, gin.H{"rewards": []interface{}{}})
}

func (h *RewardHandler) GetHistoricalINR(c *gin.Context) {
	userID := c.Param("userId")
	from := c.Query("from")
	to := c.Query("to")
	page := c.Query("page")
	size := c.Query("size")
	// TODO: Call service/repo to get historical INR values
	c.JSON(http.StatusOK, gin.H{"historical_inr": []interface{}{}})
}

func (h *RewardHandler) GetStats(c *gin.Context) {
	userID := c.Param("userId")
	// TODO: Call service/repo to get stats
	c.JSON(http.StatusOK, gin.H{"stats": gin.H{}})
}

func (h *RewardHandler) GetPortfolio(c *gin.Context) {
	userID := c.Param("userId")
	// TODO: Call service/repo to get portfolio
	c.JSON(http.StatusOK, gin.H{"portfolio": gin.H{}})
}
