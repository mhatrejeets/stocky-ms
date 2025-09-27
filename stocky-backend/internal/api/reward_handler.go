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

func (h *RewardHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/reward", h.CreateReward)
	rg.GET("/today-stocks/:userId", h.GetTodayStocks)
	rg.GET("/historical-inr/:userId", h.GetHistoricalINR)
	rg.GET("/stats/:userId", h.GetStats)
	rg.GET("/portfolio/:userId", h.GetPortfolio)
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
	rewards, err := h.Service.ListRewardsForDate(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rewards": rewards})
}

func (h *RewardHandler) GetHistoricalINR(c *gin.Context) {
	userID := c.Param("userId")
	from := c.Query("from")
	to := c.Query("to")
	page := c.Query("page")
	size := c.Query("size")
	// Set default date range if missing
	if from == "" {
		from = "2000-01-01T00:00:00Z"
	}
	if to == "" {
		to = "2100-01-01T00:00:00Z"
	}
	result, err := h.Service.GetHistoricalINR(c.Request.Context(), userID, from, to, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"historical_inr": result})
}

func (h *RewardHandler) GetStats(c *gin.Context) {
	userID := c.Param("userId")
	stats, err := h.Service.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *RewardHandler) GetPortfolio(c *gin.Context) {
	userID := c.Param("userId")
	portfolio, err := h.Service.GetPortfolio(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}
