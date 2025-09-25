package tests

import (
	"context"
	"testing"

	"github.com/mhatrejeets/stocky-ms/internal/model"
	"github.com/mhatrejeets/stocky-ms/internal/service"
	"github.com/stretchr/testify/mock"
)

func TestCreateReward_Valid(t *testing.T) {
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
	if result.Conflict {
		t.Errorf("expected no conflict, got conflict")
	}
	if result.Err != nil {
		t.Errorf("expected no error, got %v", result.Err)
	}
	if result.RewardID != "reward-uuid" {
		t.Errorf("expected reward-uuid, got %v", result.RewardID)
	}
}

func TestCreateReward_InvalidInput(t *testing.T) {
	repo := new(MockRewardRepo)
	svc := &service.RewardService{Repo: repo}
	userID := "user-1"
	req := model.CreateRewardRequest{
		StockSymbol: "RELIANCE",
		Shares:      "not-a-decimal",
		RewardedAt:  "invalid-date",
	}
	result := svc.CreateReward(context.Background(), userID, req, "")
	if result.Err == nil {
		t.Errorf("expected error for invalid input, got nil")
	}
	if result.Conflict {
		t.Errorf("expected no conflict for invalid input")
	}
}

func (m *MockRewardRepo) ListRewardsForDate(ctx context.Context, userID string, date interface{}) ([]model.Reward, error) {
	args := m.Called(ctx, userID, date)
	return args.Get(0).([]model.Reward), args.Error(1)
}

func (m *MockRewardRepo) GetHistoricalINR(ctx context.Context, userID, from, to, page, size string) ([]model.HistoricalINR, error) {
	args := m.Called(ctx, userID, from, to, page, size)
	return args.Get(0).([]model.HistoricalINR), args.Error(1)
}

func (m *MockRewardRepo) GetStats(ctx context.Context, userID string) (model.Stats, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(model.Stats), args.Error(1)
}

func (m *MockRewardRepo) GetPortfolio(ctx context.Context, userID string) (model.Portfolio, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(model.Portfolio), args.Error(1)
}
