# Currency Exchange — монорепозиторий микросервисов

Монорепозиторий из трёх сервисов:
- api_gateway — HTTP API (Gin), авторизация (JWT), CORS, rate limiting; проксирование в gRPC сервисы.
- exchange_service — gRPC сервис курсов валют, отдельный HTTP эндпоинт метрик Prometheus.
- user_service — gRPC сервис аутентификации (SSO), PostgreSQL хранилище пользователей, метрики.

Инфраструктура (локально через docker-compose):
- infra/redis — Redis
- infra/prometheus — Prometheus
- infra/grafana — Grafana
- infra/logs — Graylog/Filebeat 
- user_service/storage — PostgreSQL

## Требования
- Go 1.21+
- Docker + Docker Compose
- GNU Make 

## Быстрый старт
1) Создайте .env файлы (шаблоны ниже).
2) Поднимите инфраструктуру:
```bash
make infra-up
```
3) Запустите сервисы по отдельности в отдельных терминалах (или вкладках):
```bash
make run-api-gateway
make run-exchange
make run-user
```
4) Остановить инфраструктуру:
```bash
make infra-down
```

## Переменные окружения

### api_gateway/config/.env
```env
AUTH_PORT=50051
EXCHANGE_PORT=50052
GATEWAY_PORT=8080
TIMEOUT=5s
SECRET=super-secret
```

### exchange_service/config/.env
```env
API_KEY=your-external-api-key
PORT=50052
TIMEOUT=5s
```

### user_service/config/.env
```env
env=local
storage_path=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
tocken_ttl=24h
secret=super-secret
port=50051
timeout=5s
```

PostgreSQL поднимается из `user_service/storage/docker-compose.yaml`. Миграции лежат в `user_service/internal/migrations`.

## Основные Make-команды
- make infra-up — поднять Prometheus, Grafana, Redis, Postgres
- make infra-down — остановить инфраструктуру
- make run-api-gateway — запустить HTTP шлюз
- make run-exchange — запустить сервис обмена валют
- make run-user — запустить сервис аутентификации
- make test — прогнать тесты во всех сервисах

## Точки входа
- api_gateway/cmd/main.go — HTTP сервер на Gin
- exchange_service/cmd/main.go — gRPC сервер + метрики
- user_service/cmd/sso/main.go — gRPC сервер SSO + метрики

## Наблюдаемость
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (admin/admin при первом входе)

## Тесты
```bash
make test
```
