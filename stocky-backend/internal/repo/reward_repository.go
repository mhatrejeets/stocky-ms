package repo

import (
	"context"

	"github.com/mhatrejeets/stocky-ms/internal/model"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward model.Reward) (string, error)
	ExistsByUniqueHashOrIdempotency(ctx context.Context, uniqueHash, idempotencyKey string) (bool, string)
	// Add other methods for listing, stats, etc.
	// CheckIdempotencyKey returns (exists, response) for a given idempotency key
	CheckIdempotencyKey(ctx context.Context, key string) (bool, interface{})

	// List all rewards for a user for a given date
	ListRewardsForDate(ctx context.Context, userID string, date interface{}) ([]model.Reward, error)

	// Get historical INR values for a user
	GetHistoricalINR(ctx context.Context, userID, from, to, page, size string) ([]model.HistoricalINR, error)

	// Get stats for a user
	GetStats(ctx context.Context, userID string) (model.Stats, error)

	// Get portfolio for a user
	GetPortfolio(ctx context.Context, userID string) (model.Portfolio, error)
}
