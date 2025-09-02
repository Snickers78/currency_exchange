# API Gateway

HTTP шлюз на базе Gin. Предоставляет REST эндпоинты и проксирует запросы в gRPC сервисы `user_service` и `exchange_service`.

## Возможности
- Авторизация по JWT (`SECRET`), middleware проверки токена
- CORS и Rate limiting (`BucketLimiter`)
- Маршруты:
  - /auth/* → gRPC `user_service`
  - /exchange/* → gRPC `exchange_service`

## Конфигурация
Файл `./config/.env` со следующими переменными (см. `internal/config/config.go`):
```env
AUTH_PORT=50051          # порт gRPC user_service
EXCHANGE_PORT=50052      # порт gRPC exchange_service
GATEWAY_PORT=8080        # порт HTTP шлюза
TIMEOUT=5s               # таймаут запросов к бэкендам
SECRET=super-secret      # секрет для JWT
```

## Локальный запуск
```bash
make run-api-gateway
```
или вручную:
```bash
cd api_gateway
go run ./cmd/main.go
```

## Docker
Сборка и запуск:
```bash
cd api_gateway
docker build -t api-gateway:local -f dockerfile .
docker run --rm -p 8080:8080 --env-file ./config/.env api-gateway:local
```

Учтите, что Dockerfile по умолчанию экспонирует порт 1337, а приложение слушает `GATEWAY_PORT`. При запуске контейнера пробросьте внешний порт на `GATEWAY_PORT` из env.
