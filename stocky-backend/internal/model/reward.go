package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Reward struct {
	ID             string          `json:"id"`
	UserID         string          `json:"user_id"`
	StockSymbol    string          `json:"stock_symbol"`
	Shares         decimal.Decimal `json:"shares"`
	RewardedAt     time.Time       `json:"rewarded_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UniqueHash     string          `json:"unique_hash"`
	IdempotencyKey string          `json:"idempotency_key"`
	Status         string          `json:"status"`
}

type CreateRewardRequest struct {
	StockSymbol string `json:"stock_symbol" validate:"required"`
	Shares      string `json:"shares" validate:"required,decimal"`
	RewardedAt  string `json:"rewarded_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}
