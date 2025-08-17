package currency

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"exchange_service/infra/metrics"
	"exchange_service/internal/config"
	logg "exchange_service/internal/lib/logger"
)

var (
	ErrCurrencyNotFound = errors.New("currency not found")
	ErrInternalServer   = errors.New("internal server error")
	ErrUnsupportedCode  = errors.New("currency code unsupported")
	ErrMalformedRequest = errors.New("malformed request")
	ErrInvalidKey       = errors.New("api key is invalid")
	ErrUnknownCode      = errors.New("unknown currencycode")
)

type Exchanger struct {
	config *config.Config
	logger *slog.Logger
	cache  RatesCache
}

func NewExchanger(config *config.Config, logger *slog.Logger) *Exchanger {
	return &Exchanger{
		config: config,
		logger: logger,
		cache:  NewMemoryRatesCache(24 * time.Hour),
	}
}

func (e *Exchanger) SanitizeError(err error) string {
	if err == nil {
		return ""
	}

	errorMsg := err.Error()
	re := regexp.MustCompile(`/v6/[^/]+/`)
	sanitized := re.ReplaceAllString(errorMsg, "/v6/***/")

	return sanitized
}

func (e *Exchanger) GetExchangeRate(ctx context.Context, base_currency string, target_currency string) (float64, time.Time, error) {
	metrics.RequestCount.Inc()
	start := time.Now()

	base := strings.ToUpper(base_currency)
	target := strings.ToUpper(target_currency)

	if rates, lastUpdate, err := e.cache.GetRates(ctx, base); err == nil && rates != nil {
		if rate, ok := rates[target]; ok {
			metrics.ResponseTimeSeconds.Observe(time.Since(start).Seconds())
			return rate, lastUpdate, nil
		}

	}

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", e.config.API_KEY, base)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logEntry := logg.NewExchangeLog("error", "get_exchange_rate_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithExchangeError(e.SanitizeError(err)))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrInternalServer
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logEntry := logg.NewExchangeLog("error", "get_exchange_rate_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithExchangeError(e.SanitizeError(err)))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrInternalServer
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logEntry := logg.NewExchangeLog("error", "get_exchange_rate_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithExchangeError(e.SanitizeError(err)))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrInternalServer
	}

	var apiResp exchangeRateAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		logEntry := logg.NewExchangeLog("error", "get_exchange_rate_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithExchangeError(e.SanitizeError(err)))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrInternalServer
	}

	if apiResp.Result == "error" {
		logEntry := logg.NewExchangeLog("error", "get_exchange_rate_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithExchangeError(apiResp.ErrorType))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		switch apiResp.Result {
		case "invalid-key":
			return 0, time.Time{}, ErrInvalidKey
		case "unsupported-code":
			return 0, time.Time{}, ErrUnsupportedCode
		case "malformed-request":
			return 0, time.Time{}, ErrMalformedRequest
		case "unknown-code":
			return 0, time.Time{}, ErrUnknownCode
		default:
			return 0, time.Time{}, ErrInternalServer
		}
	}

	// Сохраняем ВСЕ курсы в кэш до следующего обновления, указанного провайдером
	lastUpdate := time.Unix(apiResp.TimeLastUpdateUnix, 0)
	nextUpdate := time.Unix(apiResp.TimeNextUpdateUnix, 0)
	_ = e.cache.SetRates(ctx, base, apiResp.Rates, lastUpdate, nextUpdate)

	rate, exists := apiResp.Rates[target]
	if !exists {
		logEntry := logg.NewExchangeLog("error", "get_exchange_rate_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithExchangeError("currency not found"))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrCurrencyNotFound
	}

	metrics.ResponseTimeSeconds.Observe(time.Since(start).Seconds())

	return rate, lastUpdate, nil
}

func (e *Exchanger) Exchange(ctx context.Context, base_currency string, target_currency string, amount float64) (float64, time.Time, error) {
	metrics.RequestCount.Inc()
	start := time.Now()

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s/%f", e.config.API_KEY, strings.ToUpper(base_currency), strings.ToUpper(target_currency), amount)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logEntry := logg.NewExchangeLog("error", "exchange_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithAmount(amount), logg.WithExchangeError(e.SanitizeError(err)))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrInternalServer
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logEntry := logg.NewExchangeLog("error", "exchange_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithAmount(amount), logg.WithExchangeError(e.SanitizeError(err)))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrInternalServer
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logEntry := logg.NewExchangeLog("error", "exchange_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithAmount(amount), logg.WithExchangeError(e.SanitizeError(err)))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrInternalServer
	}

	var apiResp exchangeCurrencyAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		logEntry := logg.NewExchangeLog("error", "exchange_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithAmount(amount), logg.WithExchangeError(e.SanitizeError(err)))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		return 0, time.Time{}, ErrInternalServer
	}

	if apiResp.Result == "error" {
		logEntry := logg.NewExchangeLog("error", "exchange_failed", logg.WithBaseCurrency(base_currency), logg.WithTargetCurrency(target_currency), logg.WithAmount(amount), logg.WithExchangeError(apiResp.ErrorType))
		e.logger.Error("error", "exchange_service", logEntry)
		metrics.ErrorCount.Inc()

		switch apiResp.Result {
		case "invalid-key":
			return 0, time.Time{}, ErrInvalidKey
		case "unsupported-code":
			return 0, time.Time{}, ErrUnsupportedCode
		case "malformed-request":
			return 0, time.Time{}, ErrMalformedRequest
		case "unknown-code":
			return 0, time.Time{}, ErrUnknownCode
		default:
			return 0, time.Time{}, ErrInternalServer
		}
	}

	updateTime := time.Unix(apiResp.TimeLastUpdateUnix, 0)
	metrics.ResponseTimeSeconds.Observe(time.Since(start).Seconds())

	return apiResp.ConversionResult, updateTime, nil
}
