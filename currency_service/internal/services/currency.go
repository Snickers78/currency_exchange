package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	logger     *slog.Logger
	apiBaseURL string
}

type exchangeRateAPIResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	Terms              string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	Rates              map[string]float64 `json:"conversion_rates"`
}

func NewService(logger *slog.Logger) *Service {
	return &Service{
		logger:     logger,
		apiBaseURL: "https://v6.exchangerate-api.com/v6/YOUR_API_KEY/latest/USD",
	}
}

func (s *Service) GetExchangeRate(ctx context.Context, currency string) (float64, time.Time, error) {
	s.logger.Info("getting exchange rate", "currency", currency)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.apiBaseURL, nil)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	var apiResp exchangeRateAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to decode API response: %w", err)
	}

	rate, exists := apiResp.Rates[currency]
	if !exists {
		return 0, time.Time{}, fmt.Errorf("currency %s not found", currency)
	}

	updateTime := time.Unix(apiResp.TimeLastUpdateUnix, 0)
	s.logger.Info("got exchange rate",
		"currency", currency,
		"rate", rate,
		"update_time", updateTime)

	return rate, updateTime, nil
}
