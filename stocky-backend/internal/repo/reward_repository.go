package repo

import (
	"context"

	"github.com/mhatrejeets/stocky-ms/internal/model"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward model.Reward) (string, error)
	ExistsByUniqueHashOrIdempotency(ctx context.Context, uniqueHash, idempotencyKey string) (bool, string)
	CheckIdempotencyKey(ctx context.Context, key string) (bool, interface{})
	ListRewardsForDate(ctx context.Context, userID string, date interface{}) ([]model.Reward, error)
	GetHistoricalINR(ctx context.Context, userID, from, to, page, size string) ([]model.HistoricalINR, error)
	GetStats(ctx context.Context, userID string) (model.Stats, error)
	GetPortfolio(ctx context.Context, userID string) (model.Portfolio, error)
}
