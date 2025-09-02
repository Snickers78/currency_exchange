# User Service (SSO)

gRPC сервис аутентификации и управления пользователями. Использует PostgreSQL как хранилище.

## Конфигурация
`./config/.env` (см. `internal/config/config.go`):
```env
env=local
storage_path=postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable
tocken_ttl=24h
secret=super-secret
port=50051
timeout=5s
```

## База данных
Поднимите PostgreSQL локально:
```bash
make pg-up
```
Или вручную:
```bash
cd user_service/storage
docker compose up -d
```
Данные сохраняются в каталоге `storage/postgres-data`.

Миграции находятся в `internal/migrations`.

## Запуск
```bash
make run-user
```
или вручную:
```bash
cd user_service
go run ./cmd/sso/main.go
```

## Метрики
- Отдельный HTTP эндпоинт метрик (см. `infra/metrics/handler.go`).
- Локальный Prometheus:
```bash
cd user_service/infra/metrics/auth_prometheus
docker compose up -d
```

## Тесты
```bash
cd user_service
go test ./...
```
