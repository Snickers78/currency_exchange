package currency

import (
	"context"
	"maps"
	"sync"
	"time"
)

// RatesCache описывает кэш курсов для базовой валюты.
// Ожидается, что реализация обеспечивает TTL и безопасна для конкурентного доступа.
type RatesCache interface {
	// GetRates возвращает карту курсов для baseCurrency, время последнего обновления
	// и ошибку. Если кэш отсутствует или истёк, должна вернуться (nil, zeroTime, nil).
	GetRates(ctx context.Context, baseCurrency string) (map[string]float64, time.Time, error)

	// SetRates сохраняет карту курсов и время последнего обновления до момента nextUpdate.
	// Если nextUpdate — нулевое время, реализация должна применить дефолтный TTL (например, 24 часа).
	SetRates(ctx context.Context, baseCurrency string, rates map[string]float64, lastUpdate time.Time, nextUpdate time.Time) error
}

// MemoryRatesCache — простая in-memory реализация кэша с TTL на элемент.
// Подходит для одиночного инстанса сервиса.
type MemoryRatesCache struct {
	mu         sync.RWMutex
	ttlDefault time.Duration
	items      map[string]memoryRatesItem
}

type memoryRatesItem struct {
	rates      map[string]float64
	lastUpdate time.Time
	expiresAt  time.Time
}

// NewMemoryRatesCache создаёт кэш с заданным дефолтным TTL.
func NewMemoryRatesCache(defaultTTL time.Duration) *MemoryRatesCache {
	if defaultTTL <= 0 {
		defaultTTL = 24 * time.Hour
	}
	return &MemoryRatesCache{
		ttlDefault: defaultTTL,
		items:      make(map[string]memoryRatesItem),
	}
}

func (c *MemoryRatesCache) GetRates(_ context.Context, baseCurrency string) (map[string]float64, time.Time, error) {
	c.mu.RLock()
	item, ok := c.items[baseCurrency]
	c.mu.RUnlock()
	if !ok {
		return nil, time.Time{}, nil
	}
	now := time.Now()
	if now.After(item.expiresAt) {
		// Истёк
		c.mu.Lock()
		delete(c.items, baseCurrency)
		c.mu.Unlock()
		return nil, time.Time{}, nil
	}
	// Возвращаем копию карты, чтобы не отдавать внутреннюю ссылку
	copied := make(map[string]float64, len(item.rates))
	maps.Copy(copied, item.rates)
	return copied, item.lastUpdate, nil
}

func (c *MemoryRatesCache) SetRates(_ context.Context, baseCurrency string, rates map[string]float64, lastUpdate time.Time, nextUpdate time.Time) error {
	expiresAt := nextUpdate
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(c.ttlDefault)
	}

	c.mu.Lock()
	c.items[baseCurrency] = memoryRatesItem{
		rates:      rates,
		lastUpdate: lastUpdate,
		expiresAt:  expiresAt,
	}
	c.mu.Unlock()
	return nil
}
