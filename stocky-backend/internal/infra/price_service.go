package infra

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type PriceProvider interface {
	GetPrice(ctx context.Context, symbol string) (decimal.Decimal, time.Time, error)
}
