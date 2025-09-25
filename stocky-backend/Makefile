build:
	go build -o stocky-backend ./cmd/api

run:
	PORT=8080 go run ./cmd/api/main.go

docker-build:
	docker build -t stocky-backend -f deploy/Dockerfile .

docker-up:
	docker-compose -f deploy/docker-compose.yml up --build
