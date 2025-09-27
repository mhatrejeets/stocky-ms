package infra

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)



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
				// Update DB stock_prices table
				if db, ok := ctx.Value("db").(*sql.DB); ok && db != nil {
					_, _ = db.ExecContext(ctx, `INSERT INTO stock_prices (symbol, price, updated_at) VALUES ($1, $2, $3)
						ON CONFLICT (symbol) DO UPDATE SET price = EXCLUDED.price, updated_at = EXCLUDED.updated_at`, symbol, val.String(), updated)
				}
			}
			time.Sleep(time.Hour)
		}
	}()
}
