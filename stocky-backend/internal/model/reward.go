package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type HistoricalINR struct {
	Date     string          `json:"date"`
	INRValue decimal.Decimal `json:"inr_value"`
	IsStale  bool            `json:"is_stale"`
}

type Stats struct {
	TodayTotalBySymbol map[string]decimal.Decimal `json:"today_total_by_symbol"`
	PortfolioValueINR  decimal.Decimal            `json:"portfolio_value_inr"`
}

type Portfolio struct {
	Holdings          []Holding       `json:"holdings"`
	PortfolioTotalINR decimal.Decimal `json:"portfolio_total_inr"`
}

type Holding struct {
	Symbol        string          `json:"symbol"`
	TotalShares   decimal.Decimal `json:"total_shares"`
	CurrentPrice  decimal.Decimal `json:"current_price"`
	TotalValueINR decimal.Decimal `json:"total_value_inr"`
}

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

type RewardCreatedEvent struct {
	RewardID      string `json:"reward_id"`
	UserID        string `json:"user_id"`
	StockSymbol   string `json:"stock_symbol"`
	Shares        string `json:"shares"`
	RewardedAt    string `json:"rewarded_at"`
	CorrelationID string `json:"correlation_id"`
}
