# stocky-backend

## Assumptions & Tool Choices

- **Language:** Go 1.21+
- **Web framework:** github.com/gin-gonic/gin
- **Logger:** github.com/sirupsen/logrus
- **DB driver:** github.com/jackc/pgx/v5/pgxpool (Postgres)
- **Decimal math:** github.com/shopspring/decimal (avoid floats for money/shares)
- **JWT:** github.com/golang-jwt/jwt/v5
- **Redis client:** github.com/redis/go-redis/v9
- **Kafka client (optional/bonus):** github.com/Shopify/sarama or confluent-kafka-go
- **Migrations:** golang-migrate/migrate (or raw SQL files)
- **Swagger/OpenAPI:** github.com/swaggo/swag + github.com/swaggo/gin-swagger
- **Metrics:** github.com/prometheus/client_golang/prometheus/promhttp
- **Tracing:** go.opentelemetry.io/otel (skeleton)
- **Rate limiting:** simple Redis-backed sliding-window or token-bucket (library or custom)
- **Testing:** testcontainers-go for integration tests (docker-compose replacement) or docker-compose for local dev.

## Folder Structure

```
stocky-backend/
├─ cmd/
│  └─ api/                # main entry
│     └─ main.go
├─ internal/
│  ├─ api/                # route handlers
│  ├─ auth/               # jwt + rbac
│  ├─ middleware/         # correlation id, logging, rate-limit, idempotency
│  ├─ service/            # business logic
│  ├─ repo/               # db queries (repository)
│  ├─ model/              # domain models / dto
│  ├─ infra/              # redis, kafka, price-service
│  └─ util/               # helpers (time, hashing)
├─ migrations/            # sql migrations
├─ docs/                  # openapi / swagger
├─ deploy/
│  ├─ docker-compose.yml
│  └─ Dockerfile
├─ scripts/
│  └─ run_migrations.sh
├─ tests/                 # integration tests
├─ .env.sample
├─ go.mod
├─ README.md
└─ Makefile
```

## Environment Variables
See `.env.sample` for required variables. **Do not commit secrets.**
