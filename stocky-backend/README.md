
# 📈 Stocky Backend

A robust, production-ready backend for managing stock rewards, portfolio tracking, and financial analytics. Built with Go, PostgreSQL, Redis, Kafka, and Gin, this system is designed for reliability, scalability, and extensibility.

---

## 🚀 Overview

Stocky Backend is a microservice that powers a stock rewards platform. It enables:

- Secure reward creation and idempotency
- Real-time portfolio and stats computation
- Historical INR value tracking
- Event-driven integrations (Kafka)
- Robust error handling and observability

---

## 🗂️ Directory Structure

```
stocky-backend/
├── cmd/
│   └── api/                # Main entrypoint (main.go)
├── internal/
│   ├── api/                # HTTP route handlers
│   ├── auth/               # JWT middleware
│   ├── infra/              # DB, Redis, Kafka, price provider
│   ├── middleware/         # Logging, rate-limit, idempotency, correlation
│   ├── model/              # Domain models & DTOs
│   ├── repo/               # Repository (DB queries)
│   └── service/            # Business logic
├── scripts/                # DB migrations, utilities
├── tests/                  # Unit & integration tests
├── deploy/                 # Docker, Compose, deployment scripts
├── .github/                # CI/CD workflows
├── .env.example            # Environment variable template
├── Makefile                # Build/test helpers
├── README.md
└── postman_collection.json # API examples
```

---

## ⚙️ Tech Stack

- **Go 1.21+**
- **Gin** (HTTP API)
- **PostgreSQL** (primary DB)
- **Redis** (idempotency, caching)
- **Kafka** (event streaming)
- **Logrus** (structured logging)
- **ShopSpring/Decimal** (precise money math)
- **Docker & Compose** (local dev, CI)
- **JWT** (auth)
- **Testcontainers** (integration tests)

---

## 🏗️ Database Schema

**Rewards Table**
- `id` (UUID, PK)
- `user_id` (string)
- `stock_symbol` (string)
- `shares` (decimal)
- `rewarded_at` (timestamp)
- `unique_hash` (string, unique)
- `idempotency_key` (string, unique)
- `status` (string)

**Stock Prices Table**
- `symbol` (string, PK)
- `price` (decimal)
- `updated_at` (timestamp)

**Ledger Entries Table**
- `id` (UUID, PK)
- `event_type` (string: reward, fee, adjustment, etc.)
- `user_id` (string)
- `stock_symbol` (string)
- `shares` (decimal)
- `inr_amount` (decimal)
- `fee_type` (string)
- `created_at` (timestamp)

**Relationships:**
- Rewards and ledger entries are linked by `user_id` and `stock_symbol`.
- Stock prices are referenced for INR calculations.

---

## 📑 API Specifications

### Authentication

All endpoints require a valid JWT in the `Authorization` header.

### Reward Creation

**POST** `/api/v1/reward`

**Request:**
```json
{
	"stock_symbol": "RELIANCE",
	"shares": "1.000000",
	"rewarded_at": "2025-09-25T11:30:00Z"
}
```
Headers:
- `Authorization: Bearer <token>`
- `X-User-ID: <user_id>`
- `Idempotency-Key: <unique-key>`

**Response:**
- `201 Created` `{ "status": "success", "reward_id": "<uuid>" }`
- `409 Conflict` `{ "error": "duplicate reward" }`

---

### Portfolio

**GET** `/api/v1/portfolio/:userId`

**Response:**
```json
{
	"holdings": [
		{
			"symbol": "RELIANCE",
			"total_shares": "2.000000",
			"current_price": "2500.00",
			"total_value_inr": "5000.00"
		}
	],
	"portfolio_total_inr": "5000.00"
}
```

---

### Stats

**GET** `/api/v1/stats/:userId`

**Response:**
```json
{
	"today_total_by_symbol": {
		"RELIANCE": "2.000000"
	},
	"portfolio_value_inr": "5000.00"
}
```

---

### Historical INR

**GET** `/api/v1/historical-inr/:userId?from=2025-09-01&to=2025-09-30`

**Response:**
```json
{
	"historical_inr": [
		{
			"date": "2025-09-25",
			"inr_value": "2500.00",
			"is_stale": false
		}
	]
}
```

---

##  System Flow

1. **Reward Creation:**  
	 - Validates input, checks idempotency (Redis + DB).
	 - Inserts reward and ledger entries.
	 - Publishes event to Kafka.
	 - Returns reward ID or conflict.

2. **Portfolio/Stats:**  
	 - Aggregates shares per symbol for the user.
	 - Fetches current prices from `stock_prices`.
	 - Computes INR values using precise decimal math.

3. **Historical INR:**  
	 - For each day, sums shares per symbol.
	 - Multiplies by current price (not historical price).
	 - Returns daily INR value.

4. **Ledger:**  
	 - Tracks all reward, fee, and adjustment events for auditability.

---

## 🛡️ Edge Case Handling

- **Duplicate/replay:**  
	- Idempotency key and unique hash prevent double-inserts.
- **Rounding errors:**  
	- All INR/share math uses `decimal.Decimal` for precision.
- **Price API downtime/stale data:**  
	- If price missing, INR value is zero; `is_stale` flag can be extended.
- **Stock splits/mergers/delisting:**  
	- Not handled in current logic (extendable via ledger adjustments).
- **Adjustments/refunds:**  
	- Not implemented, but ledger supports negative/adjustment entries.

---

## 📈 Scaling & Reliability

- **Stateless API:**  
	- All state in DB/Redis/Kafka; easy to scale horizontally.
- **Idempotency:**  
	- Safe for retries and at-least-once delivery.
- **Observability:**  
	- Structured logging, correlation IDs, and (optionally) metrics/tracing.
- **Extensibility:**  
	- Modular design for new event types, price providers, or reward logic.
- **CI/CD:**  
	- GitHub Actions for build/test; Docker for reproducible environments.

---

## 🧪 Testing

- **Unit tests:**  
	- Cover business logic and edge cases.
- **Integration tests:**  
	- Use Docker Compose and testcontainers for DB/Kafka/Redis.

---

## 📝 Deliverables

- 📦 **Public GitHub repo** with full codebase
- 📑 **API specifications** (see above)
- 🗄️ **Database schema** (see above)
- 🧠 **System flow and edge case explanations** (see above)

---

## 🛠️ Quickstart

```bash
# 1. Clone and configure .env
cp .env.example .env

# 2. Start services
docker-compose up --build

# 3. Run migrations
./scripts/run_migrations.sh

# 4. Run tests
make test

# 5. Access API at http://localhost:8080
```

---

## Note
Populate the stock_prices table with this query earlier to get correct the actuals values for stock prices and not 0.

`INSERT INTO stock_prices (symbol, price, updated_at) VALUES ('TCS','3500.00',NOW()),('INFY','1500.00',NOW()),('RELIANCE','2500.00',NOW()),('HDFCBANK','1650.50',NOW()),('ICICIBANK','1020.75',NOW()),('SBIN','650.25',NOW()),('WIPRO','420.00',NOW()),('BAJFINANCE','7200.00',NOW()),('ADANIENT','2700.00',NOW()),('HINDUNILVR','2450.00',NOW()),('ITC','455.75',NOW()),('LT','3750.00',NOW()),('AXISBANK','1125.00',NOW()),('TITAN','3400.00',NOW()),('MARUTI','10400.00',NOW()),('ONGC','195.50',NOW()),('COALINDIA','320.00',NOW()),('POWERGRID','255.30',NOW()),('JSWSTEEL','850.00',NOW()),('ASIANPAINT','3300.00',NOW()),('APPLE','195.00',NOW()),('GOOGL','135.00',NOW()),('AMZN','140.50',NOW()),('TSLA','245.00',NOW()),('MSFT','320.00',NOW()) ON CONFLICT (symbol) DO UPDATE SET price = EXCLUDED.price, updated_at = NOW();`

