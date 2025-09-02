# Exchange Service

gRPC сервис для получения курсов валют. Экспортирует метрики Prometheus отдельным HTTP сервером.

## Конфигурация
`./config/.env` (см. `internal/config/config.go`):
```env
API_KEY=your-external-api-key   # ключ к внешнему источнику курсов
PORT=50052                      # порт gRPC сервера
TIMEOUT=5s                      # таймаут исходящих запросов
```

## Запуск
```bash
make run-exchange
```
или вручную:
```bash
cd exchange_service
go run ./cmd/main.go
```

## Метрики
- HTTP эндпоинт метрик: `http://localhost:9200/metrics`
- Локальный Prometheus для сервиса можно поднять:
```bash
cd exchange_service/infra/metrics/currency_exchange_prometheus
docker compose up -d
```

## Тесты
```bash
cd exchange_service
go test ./...
```
