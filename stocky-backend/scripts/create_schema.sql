-- Rewards Table
CREATE TABLE IF NOT EXISTS rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(64) NOT NULL,
    stock_symbol VARCHAR(16) NOT NULL,
    shares NUMERIC(18,6) NOT NULL,
    rewarded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    unique_hash VARCHAR(64) NOT NULL,
    idempotency_key VARCHAR(64),
    status VARCHAR(16) NOT NULL,
    CONSTRAINT unique_reward UNIQUE (unique_hash),
    CONSTRAINT unique_idempotency UNIQUE (idempotency_key)
);

-- Stock Prices Table
CREATE TABLE IF NOT EXISTS stock_prices (
    symbol VARCHAR(16) PRIMARY KEY,
    price NUMERIC(18,4) NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Ledger Table
CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(32) NOT NULL, -- reward, fee, adjustment, etc.
    user_id VARCHAR(64),
    stock_symbol VARCHAR(16),
    shares NUMERIC(18,6),
    inr_amount NUMERIC(18,4),
    fee_type VARCHAR(32), -- brokerage, STT, GST, etc.
    created_at TIMESTAMP DEFAULT now()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_rewards_user_date ON rewards (user_id, rewarded_at);
CREATE INDEX IF NOT EXISTS idx_ledger_user ON ledger_entries (user_id);
