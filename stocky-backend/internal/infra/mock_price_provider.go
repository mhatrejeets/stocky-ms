package infra

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

var knownSymbols = []string{"RELIANCE", "TCS", "INFY"}

// MockPriceProvider returns random prices and caches in Redis

type MockPriceProvider struct {
	Redis *redis.Client
}

func (m *MockPriceProvider) GetPrice(ctx context.Context, symbol string) (decimal.Decimal, time.Time, error) {
	price, updatedAt, err := GetCachedPrice(ctx, m.Redis, symbol)
	if err == nil {
		return price, updatedAt, nil
	}
	// Generate random price
	val := decimal.NewFromFloat(rand.Float64()*1000 + 100)
	updated := time.Now()
	// Cache in Redis
	m.Redis.Set(ctx, "price:"+symbol, val.String()+","+updated.Format(time.RFC3339), 2*time.Hour)
	return val, updated, nil
}

func GetCachedPrice(ctx context.Context, rdb *redis.Client, symbol string) (decimal.Decimal, time.Time, error) {
	res, err := rdb.Get(ctx, "price:"+symbol).Result()
	if err != nil {
		return decimal.Zero, time.Time{}, err
	}
	parts := strings.Split(res, ",")
	if len(parts) != 2 {
		return decimal.Zero, time.Time{}, err
	}
	price, err := decimal.NewFromString(parts[0])
	if err != nil {
		return decimal.Zero, time.Time{}, err
	}
	updatedAt, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return decimal.Zero, time.Time{}, err
	}
	return price, updatedAt, nil
}

func StartHourlyPriceUpdater(ctx context.Context, rdb *redis.Client, symbols []string) {
	go func() {
		for {
			for _, symbol := range symbols {
				val := decimal.NewFromFloat(rand.Float64()*1000 + 100)
				updated := time.Now()
				rdb.Set(ctx, "price:"+symbol, val.String()+","+updated.Format(time.RFC3339), 2*time.Hour)
				// TODO: Update DB stock_prices table
			}
			time.Sleep(time.Hour)
		}
	}()
}
