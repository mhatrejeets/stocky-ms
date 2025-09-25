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

	// Dedup check
	exists, existingID := s.Repo.ExistsByUniqueHashOrIdempotency(ctx, uniqueHashHex, idempotencyKey)
	if exists {
		return CreateRewardResult{existingID, true, errors.New("duplicate reward")}
	}

	reward := model.Reward{
		UserID:         userID,
		StockSymbol:    req.StockSymbol,
		Shares:         shares,
		RewardedAt:     rewardedAt,
		UniqueHash:     uniqueHashHex,
		IdempotencyKey: idempotencyKey,
		Status:         "active",
	}
	id, err := s.Repo.CreateReward(ctx, reward)
	if err != nil {
		return CreateRewardResult{"", false, err}
	}
	return CreateRewardResult{id, false, nil}
}
