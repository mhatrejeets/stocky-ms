package repo

import (
	"context"
	"stocky-backend/internal/model"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, reward model.Reward) (string, error)
	ExistsByUniqueHashOrIdempotency(ctx context.Context, uniqueHash, idempotencyKey string) (bool, string)
	// Add other methods for listing, stats, etc.
}
