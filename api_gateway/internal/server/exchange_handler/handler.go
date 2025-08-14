package handler

import (
	"api_gateway/infra/kafka"
	"api_gateway/internal/exchange"
	exchangev1 "api_gateway/internal/gen/exchange/proto"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Exchange struct {
	exchangeClient *exchange.ExchangeClient
	hook           *kafka.KafkaHook
}

func NewExchangeHandler(router *gin.Engine, exchangeClient *exchange.ExchangeClient, hook *kafka.KafkaHook) {
	handler := &Exchange{
		exchangeClient: exchangeClient,
		hook:           hook,
	}
	router.POST("/exchange/rate", handler.GetExchangeRate)
	router.POST("/exchange", handler.Exchange)
}

func (e *Exchange) GetExchangeRate(c *gin.Context) {
	var req rateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := NewExchangeLog("error", "rate_invalid_request", WithBaseCurrency(req.BaseCurrency), WithTargetCurrency(req.TargetCurrency), WithExchangeError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			e.hook.Fire(string(msg))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := e.exchangeClient.GetExchangeRate(c.Request.Context(), &exchangev1.ExchangeRateRequest{
		BaseCurrency:   req.BaseCurrency,
		TargetCurrency: req.TargetCurrency,
	})
	if err != nil {
		logEntry := NewExchangeLog("error", "rate_grpc_failed", WithBaseCurrency(req.BaseCurrency), WithTargetCurrency(req.TargetCurrency), WithExchangeError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			e.hook.Fire(string(msg))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get rate failed"})
		return
	}

	logEntry := NewExchangeLog("info", "rate_success", WithBaseCurrency(req.BaseCurrency), WithTargetCurrency(req.TargetCurrency), WithRate(resp.Rate), WithCurrencyName(resp.CurrencyName))
	if msg, err := json.Marshal(logEntry); err == nil {
		e.hook.Fire(string(msg))
	}

	c.JSON(http.StatusOK, gin.H{
		"currency_name": resp.CurrencyName,
		"rate":          resp.Rate,
		"timestamp":     resp.Timestamp,
	})
}

func (e *Exchange) Exchange(c *gin.Context) {
	var req exchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := NewExchangeLog("error", "exchange_invalid_request", WithBaseCurrency(req.BaseCurrency), WithTargetCurrency(req.TargetCurrency), WithAmount(req.Amount), WithExchangeError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			e.hook.Fire(string(msg))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := e.exchangeClient.Exchange(c.Request.Context(), &exchangev1.ExchangeRequest{
		BaseCurrency:   req.BaseCurrency,
		TargetCurrency: req.TargetCurrency,
		Amount:         req.Amount,
	})
	if err != nil {
		logEntry := NewExchangeLog("error", "exchange_grpc_failed", WithBaseCurrency(req.BaseCurrency), WithTargetCurrency(req.TargetCurrency), WithAmount(req.Amount), WithExchangeError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			e.hook.Fire(string(msg))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "exchange failed"})
		return
	}

	logEntry := NewExchangeLog("info", "exchange_success", WithBaseCurrency(req.BaseCurrency), WithTargetCurrency(req.TargetCurrency), WithAmount(resp.Amount), WithCurrencyName(resp.Currency))
	if msg, err := json.Marshal(logEntry); err == nil {
		e.hook.Fire(string(msg))
	}

	c.JSON(http.StatusOK, gin.H{
		"currency":  resp.Currency,
		"amount":    resp.Amount,
		"timestamp": resp.Timestamp,
	})
}
