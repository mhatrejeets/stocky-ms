package main

import (
	"os"


	"stocky-backend/internal/repo"
	"stocky-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logrus global logger
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Wire up RewardHandler (using dummy repo for now)
	var dummyRepo repo.RewardRepository // TODO: Replace with actual implementation
	rewardService := &service.RewardService{Repo: dummyRepo}
	rewardHandler := &api.RewardHandler{Service: rewardService}
	rewardHandler.RegisterRoutes(r)

	logrus.Infof("Starting server on port %s", port)
	r.Run(":" + port)
}
