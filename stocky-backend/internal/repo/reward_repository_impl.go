package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"

	"github.com/mhatrejeets/stocky-ms/internal/model"
	"github.com/sirupsen/logrus"
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

// GetPortfolio returns the user's portfolio: total shares per stock and their current value
func (r *RewardRepositoryImpl) GetPortfolio(ctx context.Context, userID string) (model.Portfolio, error) {
	query := `SELECT stock_symbol, SUM(shares) as total_shares FROM rewards WHERE user_id = $1 GROUP BY stock_symbol`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return model.Portfolio{}, err
	}
	defer rows.Close()
	var holdings []model.Holding
	var totalValue decimal.Decimal
	for rows.Next() {
		var symbol string
		var sharesStr string
		err := rows.Scan(&symbol, &sharesStr)
		if err != nil {
			return model.Portfolio{}, err
		}
		shares, _ := decimal.NewFromString(sharesStr)
		// Fetch current price from stock_prices table
		var priceStr string
		err = r.DB.QueryRowContext(ctx, "SELECT price FROM stock_prices WHERE symbol = $1", symbol).Scan(&priceStr)
		var price decimal.Decimal
		if err == nil {
			price, _ = decimal.NewFromString(priceStr)
		} else {
			price = decimal.Zero
		}
		value := shares.Mul(price)
		totalValue = totalValue.Add(value)
		holdings = append(holdings, model.Holding{
			Symbol:        symbol,
			TotalShares:   shares,
			CurrentPrice:  price,
			TotalValueINR: value,
		})
	}
	portfolio := model.Portfolio{
		Holdings:          holdings,
		PortfolioTotalINR: totalValue,
	}
	return portfolio, nil
}

func (r *RewardRepositoryImpl) CreateReward(ctx context.Context, reward model.Reward) (string, error) {
	// Insert into rewards table
	query := `INSERT INTO rewards (
	       id, user_id, stock_symbol, shares, rewarded_at, created_at, unique_hash, idempotency_key, status
       ) VALUES (
	       $1, $2, $3, $4, $5, $6, $7, $8, $9
       ) RETURNING id`
	var id string
	err := r.DB.QueryRowContext(ctx, query,
		reward.ID,
		reward.UserID,
		reward.StockSymbol,
		reward.Shares.String(),
		reward.RewardedAt,
		reward.CreatedAt,
		reward.UniqueHash,
		reward.IdempotencyKey,
		reward.Status,
	).Scan(&id)
	if err != nil {
		logrus.WithError(err).Error("Failed to insert reward")
		return "", err
	}

	// Insert ledger entries for reward
	// Example: record stock units, INR outflow, and company fees
	ledgerQuery := `INSERT INTO ledger_entries (
	       event_type, user_id, stock_symbol, shares, inr_amount, fee_type, created_at
       ) VALUES (
	       $1, $2, $3, $4, $5, $6, $7
       )`
	// Mock price provider: get current price
	price := decimal.NewFromFloat(3000.00) // Replace with real price lookup
	inrAmount := price.Mul(reward.Shares)
	// Record stock purchase
	if _, err := r.DB.ExecContext(ctx, ledgerQuery,
		"reward", reward.UserID, reward.StockSymbol, reward.Shares.String(), inrAmount.String(), "", reward.CreatedAt); err != nil {
		logrus.WithError(err).Error("Failed to insert ledger entry: reward purchase")
		return "", err
	}
	// Record brokerage fee (example: 0.1%)
	brokerage := inrAmount.Mul(decimal.NewFromFloat(0.001))
	if _, err := r.DB.ExecContext(ctx, ledgerQuery,
		"fee", reward.UserID, reward.StockSymbol, reward.Shares.String(), brokerage.String(), "brokerage", reward.CreatedAt); err != nil {
		logrus.WithError(err).Error("Failed to insert ledger entry: brokerage fee")
		return "", err
	}
	// Record STT fee (example: 0.025%)
	stt := inrAmount.Mul(decimal.NewFromFloat(0.00025))
	if _, err := r.DB.ExecContext(ctx, ledgerQuery,
		"fee", reward.UserID, reward.StockSymbol, reward.Shares.String(), stt.String(), "STT", reward.CreatedAt); err != nil {
		logrus.WithError(err).Error("Failed to insert ledger entry: STT fee")
		return "", err
	}

	// Set idempotency key in Redis (if needed)

	// Publish to Kafka
	if r.KafkaProducer != nil {
		event := model.RewardCreatedEvent{
			RewardID:      id,
			UserID:        reward.UserID,
			StockSymbol:   reward.StockSymbol,
			Shares:        reward.Shares.String(),
			RewardedAt:    reward.RewardedAt.Format(time.RFC3339),
			CorrelationID: reward.IdempotencyKey,
		}
		_ = r.KafkaProducer.PublishRewardCreated(ctx, event)
	}
	return id, nil
}

func (r *RewardRepositoryImpl) ExistsByUniqueHashOrIdempotency(ctx context.Context, uniqueHash, idempotencyKey string) (bool, string) {
	// Check for existing reward by unique hash or idempotency key
	query := `SELECT id FROM rewards WHERE unique_hash = $1 OR idempotency_key = $2 LIMIT 1`
	var id string
	err := r.DB.QueryRowContext(ctx, query, uniqueHash, idempotencyKey).Scan(&id)
	if err == sql.ErrNoRows {
		return false, ""
	}
	if err != nil {
		logrus.WithError(err).Error("Failed to check for existing reward")
		return false, ""
	}
	return true, id
}

func (r *RewardRepositoryImpl) CheckIdempotencyKey(ctx context.Context, key string) (bool, interface{}) {
	// Check for existing idempotency key
	query := `SELECT id FROM rewards WHERE idempotency_key = $1 LIMIT 1`
	var id string
	err := r.DB.QueryRowContext(ctx, query, key).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		logrus.WithError(err).Error("Failed to check idempotency key")
		return false, nil
	}
	return true, id
}

func (r *RewardRepositoryImpl) ListRewardsForDate(ctx context.Context, userID string, date interface{}) ([]model.Reward, error) {
	// Query rewards for the given user and date
	t, ok := date.(time.Time)
	if !ok {
		return nil, nil
	}
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := start.Add(24 * time.Hour)
	query := `SELECT id, user_id, stock_symbol, shares, rewarded_at, created_at, unique_hash, idempotency_key, status FROM rewards WHERE user_id = $1 AND rewarded_at >= $2 AND rewarded_at < $3`
	rows, err := r.DB.QueryContext(ctx, query, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rewards []model.Reward
	for rows.Next() {
		var rw model.Reward
		var sharesStr string
		err := rows.Scan(&rw.ID, &rw.UserID, &rw.StockSymbol, &sharesStr, &rw.RewardedAt, &rw.CreatedAt, &rw.UniqueHash, &rw.IdempotencyKey, &rw.Status)
		if err != nil {
			return nil, err
		}
		rw.Shares, _ = decimal.NewFromString(sharesStr)
		rewards = append(rewards, rw)
	}
	return rewards, nil
}

func (r *RewardRepositoryImpl) GetHistoricalINR(ctx context.Context, userID, from, to, page, size string) ([]model.HistoricalINR, error) {
	// Query historical INR values for a user
	// For each day, sum shares per symbol, then multiply by current price
	query := `SELECT to_char(rewarded_at, 'YYYY-MM-DD') as date, stock_symbol, SUM(shares) as total_shares FROM rewards WHERE user_id = $1 AND rewarded_at >= $2 AND rewarded_at <= $3 GROUP BY date, stock_symbol ORDER BY date`
	rows, err := r.DB.QueryContext(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Map date -> symbol -> shares
	type daySymbol struct {
		date   string
		symbol string
	}
	dayShares := make(map[daySymbol]decimal.Decimal)
	dateSet := make(map[string]struct{})
	for rows.Next() {
		var date, symbol, sharesStr string
		err := rows.Scan(&date, &symbol, &sharesStr)
		if err != nil {
			return nil, err
		}
		shares, _ := decimal.NewFromString(sharesStr)
		dayShares[daySymbol{date, symbol}] = shares
		dateSet[date] = struct{}{}
	}
	// For each date, sum INR value using current price
	var result []model.HistoricalINR
	for date := range dateSet {
		var totalINR decimal.Decimal
		for ds, shares := range dayShares {
			if ds.date == date {
				// Get current price for symbol
				var priceStr string
				err := r.DB.QueryRowContext(ctx, "SELECT price FROM stock_prices WHERE symbol = $1", ds.symbol).Scan(&priceStr)
				var price decimal.Decimal
				if err == nil {
					price, _ = decimal.NewFromString(priceStr)
				} else {
					price = decimal.Zero
				}
				totalINR = totalINR.Add(shares.Mul(price))
			}
		}
		result = append(result, model.HistoricalINR{
			Date:     date,
			INRValue: totalINR,
			IsStale:  false,
		})
	}
	// Sort result by date ascending
	// (optional, for consistent output)
	// sort.Slice(result, func(i, j int) bool { return result[i].Date < result[j].Date })
	return result, nil
}

func (r *RewardRepositoryImpl) GetStats(ctx context.Context, userID string) (model.Stats, error) {
	query := `SELECT stock_symbol, SUM(shares) as total_shares FROM rewards WHERE user_id = $1 GROUP BY stock_symbol`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return model.Stats{}, err
	}
	defer rows.Close()
	stats := model.Stats{TodayTotalBySymbol: make(map[string]decimal.Decimal)}
	var total decimal.Decimal
	for rows.Next() {
		var symbol string
		var sharesStr string
		err := rows.Scan(&symbol, &sharesStr)
		if err != nil {
			return model.Stats{}, err
		}
		shares, _ := decimal.NewFromString(sharesStr)
		stats.TodayTotalBySymbol[symbol] = shares
		// Fetch current price from stock_prices table
		var priceStr string
		err = r.DB.QueryRowContext(ctx, "SELECT price FROM stock_prices WHERE symbol = $1", symbol).Scan(&priceStr)
		var price decimal.Decimal
		if err == nil {
			price, _ = decimal.NewFromString(priceStr)
		} else {
			price = decimal.Zero
		}
		total = total.Add(shares.Mul(price))
	}
	stats.PortfolioValueINR = total
	return stats, nil
}
