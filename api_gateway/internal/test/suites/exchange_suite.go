package test

import (
	"api_gateway/internal/config"
	"api_gateway/internal/exchange"
	"context"
	"testing"
	"time"
)

const (
	Test = true
)

type ExchangeSuite struct {
	Test   *testing.T
	Client *exchange.ExchangeClient
	Cfg    *config.Config
}

func NewExchangeSuite(t *testing.T) (context.Context, *ExchangeSuite) {
	t.Helper()
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cfg := config.LoadConfig(Test)

	client := exchange.NewExchangeClient(cfg)

	return ctx, &ExchangeSuite{
		Test:   t,
		Client: client,
		Cfg:    cfg,
	}
}
