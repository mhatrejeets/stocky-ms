package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/mhatrejeets/stocky-ms/internal/model"
)

type RewardRepositoryImpl struct {
	DB            *sql.DB
	Redis         RedisIdempotencyStore
	KafkaProducer KafkaPublisher // interface for Kafka
}

// KafkaPublisher interface
type KafkaPublisher interface {
	PublishRewardCreated(ctx context.Context, event model.RewardCreatedEvent) error
}

// RedisIdempotencyStore interface
type RedisIdempotencyStore interface {
	SetIfNotExists(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	Get(ctx context.Context, key string) (string, error)
}

func (r *RewardRepositoryImpl) CreateReward(ctx context.Context, reward model.Reward) (string, error) {
	// Insert into rewards table
	// Insert ledger entries
	// Set idempotency key in Redis
	// Publish to Kafka
	if r.KafkaProducer != nil {
		event := model.RewardCreatedEvent{
			RewardID:      reward.ID,
			UserID:        reward.UserID,
			StockSymbol:   reward.StockSymbol,
			Shares:        reward.Shares.String(),
			RewardedAt:    reward.RewardedAt.Format(time.RFC3339),
			CorrelationID: reward.IdempotencyKey,
		}
		_ = r.KafkaProducer.PublishRewardCreated(ctx, event)
	}
	return "reward-id", nil
}

func (r *RewardRepositoryImpl) ExistsByUniqueHashOrIdempotency(ctx context.Context, uniqueHash, idempotencyKey string) (bool, string) {
	return false, ""
}

func (r *RewardRepositoryImpl) CheckIdempotencyKey(ctx context.Context, key string) (bool, interface{}) {
	return false, nil
}

func (r *RewardRepositoryImpl) ListRewardsForDate(ctx context.Context, userID string, date interface{}) ([]model.Reward, error) {
	return []model.Reward{}, nil
}

func (r *RewardRepositoryImpl) GetHistoricalINR(ctx context.Context, userID, from, to, page, size string) ([]model.HistoricalINR, error) {
	return []model.HistoricalINR{}, nil
}

func (r *RewardRepositoryImpl) GetStats(ctx context.Context, userID string) (model.Stats, error) {
	return model.Stats{}, nil
}

func (r *RewardRepositoryImpl) GetPortfolio(ctx context.Context, userID string) (model.Portfolio, error) {
	return model.Portfolio{}, nil
}
