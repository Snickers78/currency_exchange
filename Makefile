SHELL := /bin/sh

# -------- Infrastructure --------
.PHONY: infra-up infra-down redis-up redis-down prom-up prom-down grafana-up grafana-down pg-up pg-down

infra-up: redis-up prom-up grafana-up pg-up

infra-down: redis-down prom-down grafana-down pg-down

redis-up:
	docker compose -f infra/redis/docker-compose.yaml up -d

redis-down:
	docker compose -f infra/redis/docker-compose.yaml down

prom-up:
	docker compose -f infra/prometheus/docker-compose.yaml up -d

prom-down:
	docker compose -f infra/prometheus/docker-compose.yaml down

grafana-up:
	docker compose -f infra/grafana/docker-compose.yaml up -d

grafana-down:
	docker compose -f infra/grafana/docker-compose.yaml down 

pg-up:
	docker compose -f user_service/storage/docker-compose.yaml --env-file user_service/config/.env up -d

pg-down:
	docker compose -f user_service/storage/docker-compose.yaml down 

# -------- Services --------
.PHONY: run-api-gateway run-exchange run-user

run-api-gateway:
	cd api_gateway && go run ./cmd/main.go

run-exchange:
	cd exchange_service && go run ./cmd/main.go

run-user:
	cd user_service && go run ./cmd/sso/main.go

