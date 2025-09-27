package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/mhatrejeets/stocky-ms/internal/model"
	"github.com/mhatrejeets/stocky-ms/internal/repo"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var validate = validator.New()

// RewardService provides business logic for rewards

//go:generate mockgen -source=reward_service.go -destination=../mocks/mock_reward_service.go -package=mocks

type RewardService struct {
	Repo repo.RewardRepository
}

type CreateRewardResult struct {
	RewardID string
	Conflict bool
	Err      error
}

func (s *RewardService) CreateReward(ctx context.Context, userID string, req model.CreateRewardRequest, idempotencyKey string) CreateRewardResult {
	if err := validate.Struct(req); err != nil {
		return CreateRewardResult{"", false, err}
	}
	shares, err := decimal.NewFromString(req.Shares)
	if err != nil {
		return CreateRewardResult{"", false, errors.New("invalid shares format")}
	}
	rewardedAt, err := time.Parse(time.RFC3339, req.RewardedAt)
	if err != nil {
		return CreateRewardResult{"", false, errors.New("invalid rewarded_at format")}
	}
	// Compute unique hash
	uniqueStr := userID + req.StockSymbol + req.Shares + req.RewardedAt
	uniqueHash := sha256.Sum256([]byte(uniqueStr))
	uniqueHashHex := hex.EncodeToString(uniqueHash[:])

	exists, existingID := s.Repo.ExistsByUniqueHashOrIdempotency(ctx, uniqueHashHex, idempotencyKey)
	if exists {
		return CreateRewardResult{existingID, true, errors.New("duplicate reward")}
	}

	reward := model.Reward{
		ID:             uuid.NewString(),
		UserID:         userID,
		StockSymbol:    req.StockSymbol,
		Shares:         shares,
		RewardedAt:     rewardedAt,
		UniqueHash:     uniqueHashHex,
		IdempotencyKey: idempotencyKey,
		Status:         "active",
		CreatedAt:      time.Now(),
	}
	id, err := s.Repo.CreateReward(ctx, reward)
	if err != nil {
		return CreateRewardResult{"", false, err}
	}
	return CreateRewardResult{id, false, nil}
}

// List all rewards for today for a user
func (s *RewardService) ListRewardsForDate(ctx context.Context, userID string) ([]model.Reward, error) {
	return s.Repo.ListRewardsForDate(ctx, userID, time.Now())
}

// Get historical INR values for a user
func (s *RewardService) GetHistoricalINR(ctx context.Context, userID, from, to, page, size string) ([]model.HistoricalINR, error) {
	return s.Repo.GetHistoricalINR(ctx, userID, from, to, page, size)
}

// Get stats for a user
func (s *RewardService) GetStats(ctx context.Context, userID string) (model.Stats, error) {
	return s.Repo.GetStats(ctx, userID)
}

// Get portfolio for a user
func (s *RewardService) GetPortfolio(ctx context.Context, userID string) (model.Portfolio, error) {
	return s.Repo.GetPortfolio(ctx, userID)
}
