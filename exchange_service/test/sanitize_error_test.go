package test

import (
	"errors"
	"log/slog"
	"testing"

	"exchange_service/internal/config"
	currency "exchange_service/internal/services"
)

func TestSanitizeError(t *testing.T) {
	config := &config.Config{API_KEY: "test-api-key-12345"}
	exchanger := currency.NewExchanger(config, slog.Default())

	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "nil error",
			input:    nil,
			expected: "",
		},
		{
			name:     "error without API key",
			input:    errors.New("simple error message"),
			expected: "simple error message",
		},
		{
			name:     "error with API key in URL",
			input:    errors.New("failed to fetch: https://v6.exchangerate-api.com/v6/test-api-key-12345/latest/USD"),
			expected: "failed to fetch: https://v6.exchangerate-api.com/v6/***/latest/USD",
		},
		{
			name:     "error with API key in pair URL",
			input:    errors.New("failed to fetch: https://v6.exchangerate-api.com/v6/test-api-key-12345/pair/USD/EUR/100.0"),
			expected: "failed to fetch: https://v6.exchangerate-api.com/v6/***/pair/USD/EUR/100.0",
		},
		{
			name:     "error with multiple API keys",
			input:    errors.New("multiple URLs: https://v6.exchangerate-api.com/v6/key1/latest/USD and https://v6.exchangerate-api.com/v6/key2/pair/USD/EUR/100.0"),
			expected: "multiple URLs: https://v6.exchangerate-api.com/v6/***/latest/USD and https://v6.exchangerate-api.com/v6/***/pair/USD/EUR/100.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := exchanger.SanitizeError(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeError() = %v, want %v", result, tt.expected)
			}
		})
	}
}
