package test

import (
	"api_gateway/internal/lib/logger"
	handler "api_gateway/internal/server/exchange_handler"
	test "api_gateway/internal/test/suites"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ExchangeRateResponse struct {
	BaseCurreny    string `json:"base_currency"`
	TargetCurrency string `json:"target_currency"`
	Amount         string `json:"amount, omitempty"`
	Timestamp      string `json:"timestamp"`
}

func TestHttpRequestExchangeRate(t *testing.T) {
	ctx, s := test.NewExchangeSuite(t)
	router := gin.Default()
	log := logger.InitLogger("local")
	reqBody := `{"base_currency": "USD", "target_currency": "RUB"}`

	handler.NewExchangeHandler(router, s.Client, log, s.Cfg.Secret)

	server := httptest.NewServer(router)
	defer server.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/exchange/rate", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Error starting http server: %v", err)
	}
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Inpob3BhMTIzQHlhbmRleC5ydSIsImV4cCI6MTc1NTI4MzcxOSwiaWQiOjF9.qTxECJQd98mnaWIkiMHry9BJf8GR5NWhLRWaV8o4u1o")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error getting http response: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error getting response body: %v", err)
	}

	var response ExchangeRateResponse

	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}

	assert.Empty(t, response, "Response body should be not nil")

}

func TestHttpRequestExchange(t *testing.T) {
	ctx, s := test.NewExchangeSuite(t)
	router := gin.Default()
	log := logger.InitLogger("local")
	reqBody := `{"base_currency": "USD", "target_currency": "RUB", "amount": 123}`

	handler.NewExchangeHandler(router, s.Client, log, s.Cfg.Secret)

	server := httptest.NewServer(router)
	defer server.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/exchange", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Error starting http server: %v", err)
	}
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Inpob3BhMTIzQHlhbmRleC5ydSIsImV4cCI6MTc1NTI4MzcxOSwiaWQiOjF9.qTxECJQd98mnaWIkiMHry9BJf8GR5NWhLRWaV8o4u1o")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error getting http response: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error getting response body: %v", err)
	}

	var response ExchangeRateResponse

	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Error unmarshalling json: %v", err)
	}

	assert.Empty(t, response.Amount, "Amount should be non zero")
	assert.Empty(t, response, "Response body should be non nil")

}
