package main

import (
	"os"

	"github.com/mhatrejeets/stocky-ms/internal/api"
	"github.com/mhatrejeets/stocky-ms/internal/auth"
	"github.com/mhatrejeets/stocky-ms/internal/infra"
	"github.com/mhatrejeets/stocky-ms/internal/repo"
	"github.com/mhatrejeets/stocky-ms/internal/service"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logrus global logger
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	// Initialize DB
	db, err := infra.NewDB()
	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient := infra.NewRedisClient()
	defer redisClient.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Wire up RewardHandler with real repository implementation
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	producer, err := sarama.NewSyncProducer(brokers, kafkaConfig)
	if err != nil {
		logrus.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	kafkaProducer := &infra.KafkaProducer{Producer: producer}

	// Redis idempotency implementation
	redisIdem := &infra.RedisIdempotencyStoreImpl{Client: redisClient}

	repoImpl := &repo.RewardRepositoryImpl{
		DB:            db,
		Redis:         redisIdem,
		KafkaProducer: kafkaProducer,
	}

	rewardService := &service.RewardService{Repo: repoImpl}
	rewardHandler := &api.RewardHandler{Service: rewardService}
	// Protect API routes with JWT middleware
	jwtSecret := os.Getenv("JWT_SECRET")
	v1 := r.Group("/api/v1", auth.JWT(jwtSecret))
	rewardHandler.RegisterRoutes(v1)

	logrus.Infof("Starting server on port %s", port)
	r.Run(":" + port)
}
