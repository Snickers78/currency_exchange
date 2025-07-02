package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"exchange_service/internal/config"
	"exchange_service/internal/metrics"
)

type Exchanger struct {
	logger *slog.Logger
	config *config.Config
}

func NewExchanger(logger *slog.Logger, config *config.Config) *Exchanger {
	return &Exchanger{
		logger: logger,
		config: config,
	}
}

func (e *Exchanger) GetExchangeRate(ctx context.Context, base_currency string, target_currency string) (float64, time.Time, error) {
	e.logger.Info("getting exchange rate", "currency", target_currency)
	metrics.RequestCount.Inc()
	start := time.Now()
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", e.config.API_KEY, strings.ToUpper(base_currency))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		e.logger.Error("failed to create request", "error", err)
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp exchangeRateAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("failed to decode API response: %w", err)
	}

	rate, exists := apiResp.Rates[strings.ToUpper(target_currency)]
	if !exists {
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("currency %s not found", target_currency)
	}

	updateTime := time.Unix(apiResp.TimeLastUpdateUnix, 0)

	metrics.ResponseTimeSeconds.Observe(time.Since(start).Seconds())

	return rate, updateTime, nil
}

func (e *Exchanger) Exchange(ctx context.Context, base_currency string, target_currency string, amount float64) (float64, time.Time, error) {
	metrics.RequestCount.Inc()
	start := time.Now()
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s/%f", e.config.API_KEY, strings.ToUpper(base_currency), strings.ToUpper(target_currency), amount)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		e.logger.Error("failed to create request", "error", err)
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp exchangeCurrencyAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("failed to decode API response: %w", err)
	}

	conversionResult := apiResp.ConversionResult
	if conversionResult == 0 {
		metrics.ErrorCount.Inc()
		return 0, time.Time{}, fmt.Errorf("currency %s not found", target_currency)
	}

	updateTime := time.Unix(apiResp.TimeLastUpdateUnix, 0)

	metrics.ResponseTimeSeconds.Observe(time.Since(start).Seconds())

	return conversionResult, updateTime, nil
}
