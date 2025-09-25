# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o stocky-backend ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/stocky-backend .
COPY .env.sample .
EXPOSE 8080
CMD ["./stocky-backend"]
