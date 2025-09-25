package tests

import (
	"context"
	"testing"

	"github.com/mhatrejeets/stocky-ms/internal/model"
	"github.com/mhatrejeets/stocky-ms/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRewardRepo struct {
	mock.Mock
}

func (m *MockRewardRepo) CreateReward(ctx context.Context, reward model.Reward) (string, error) {
	args := m.Called(ctx, reward)
	return args.String(0), args.Error(1)
}
func (m *MockRewardRepo) ExistsByUniqueHashOrIdempotency(ctx context.Context, uniqueHash, idempotencyKey string) (bool, string) {
	args := m.Called(ctx, uniqueHash, idempotencyKey)
	return args.Bool(0), args.String(1)
}
func (m *MockRewardRepo) CheckIdempotencyKey(ctx context.Context, key string) (bool, interface{}) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Get(1)
}

func TestCreateReward_Success(t *testing.T) {
	repo := new(MockRewardRepo)
	svc := &service.RewardService{Repo: repo}
	userID := "user-1"
	req := model.CreateRewardRequest{
		StockSymbol: "RELIANCE",
		Shares:      "1.000000",
		RewardedAt:  "2025-09-25T11:30:00Z",
	}
	repo.On("ExistsByUniqueHashOrIdempotency", mock.Anything, mock.Anything, mock.Anything).Return(false, "")
	repo.On("CreateReward", mock.Anything, mock.Anything).Return("reward-uuid", nil)
	result := svc.CreateReward(context.Background(), userID, req, "")
	assert.False(t, result.Conflict)
	assert.NoError(t, result.Err)
	assert.Equal(t, "reward-uuid", result.RewardID)
}

func TestCreateReward_Duplicate(t *testing.T) {
	repo := new(MockRewardRepo)
	svc := &service.RewardService{Repo: repo}
	userID := "user-1"
	req := model.CreateRewardRequest{
		StockSymbol: "RELIANCE",
		Shares:      "1.000000",
		RewardedAt:  "2025-09-25T11:30:00Z",
	}
	repo.On("ExistsByUniqueHashOrIdempotency", mock.Anything, mock.Anything, mock.Anything).Return(true, "reward-uuid")
	result := svc.CreateReward(context.Background(), userID, req, "")
	assert.True(t, result.Conflict)
	assert.Error(t, result.Err)
	assert.Equal(t, "reward-uuid", result.RewardID)
}
