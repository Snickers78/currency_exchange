package handler

import (
	"api_gateway/internal/exchange"
	exchangev1 "api_gateway/internal/gen/exchange/proto"
	"api_gateway/internal/kafka"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Exchange struct {
	exchangeClient *exchange.ExchangeClient
	hook           *kafka.KafkaHook
}

func NewExchangeHandler(router *gin.Engine, exchangeClient *exchange.ExchangeClient, hook *kafka.KafkaHook) *Exchange {
	handler := &Exchange{
		exchangeClient: exchangeClient,
		hook:           hook,
	}
	router.POST("/exchange/rate", handler.GetExchangeRate)
	router.POST("/exchange", handler.Exchange)
	return handler
}

func (e *Exchange) GetExchangeRate(c *gin.Context) {
	var req rateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := ExchangeLog{
			Level:          "error",
			Event:          "rate_invalid_request",
			BaseCurrency:   req.BaseCurrency,
			TargetCurrency: req.TargetCurrency,
			Error:          err.Error(),
			Time:           time.Now().Format(time.RFC3339),
		}
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
		logEntry := ExchangeLog{
			Level:          "error",
			Event:          "rate_grpc_failed",
			BaseCurrency:   req.BaseCurrency,
			TargetCurrency: req.TargetCurrency,
			Error:          err.Error(),
			Time:           time.Now().Format(time.RFC3339),
		}
		if msg, err := json.Marshal(logEntry); err == nil {
			e.hook.Fire(string(msg))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get rate failed"})
		return
	}

	logEntry := ExchangeLog{
		Level:          "info",
		Event:          "rate_success",
		BaseCurrency:   req.BaseCurrency,
		TargetCurrency: req.TargetCurrency,
		Rate:           resp.Rate,
		CurrencyName:   resp.CurrencyName,
		Time:           time.Now().Format(time.RFC3339),
	}
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
		logEntry := ExchangeLog{
			Level:          "error",
			Event:          "exchange_invalid_request",
			BaseCurrency:   req.BaseCurrency,
			TargetCurrency: req.TargetCurrency,
			Amount:         req.Amount,
			Error:          err.Error(),
			Time:           time.Now().Format(time.RFC3339),
		}
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
		logEntry := ExchangeLog{
			Level:          "error",
			Event:          "exchange_grpc_failed",
			BaseCurrency:   req.BaseCurrency,
			TargetCurrency: req.TargetCurrency,
			Amount:         req.Amount,
			Error:          err.Error(),
			Time:           time.Now().Format(time.RFC3339),
		}
		if msg, err := json.Marshal(logEntry); err == nil {
			e.hook.Fire(string(msg))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "exchange failed"})
		return
	}

	logEntry := ExchangeLog{
		Level:          "info",
		Event:          "exchange_success",
		BaseCurrency:   req.BaseCurrency,
		TargetCurrency: req.TargetCurrency,
		Amount:         resp.Amount,
		CurrencyName:   resp.Currency,
		Time:           time.Now().Format(time.RFC3339),
	}
	if msg, err := json.Marshal(logEntry); err == nil {
		e.hook.Fire(string(msg))
	}

	c.JSON(http.StatusOK, gin.H{
		"currency":  resp.Currency,
		"amount":    resp.Amount,
		"timestamp": resp.Timestamp,
	})
}
