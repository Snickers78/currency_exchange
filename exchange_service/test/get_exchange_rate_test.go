package test

import (
	currency_v1 "exchange_service/gen/proto"
	"exchange_service/test/suite"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExchange_CornerCases(t *testing.T) {
	ctx, s := suite.NewSuite(t)

	t.Run("nonexistent currency", func(t *testing.T) {
		resp, err := s.Client.GetExchangeRate(ctx, &currency_v1.ExchangeRateRequest{
			BaseCurrency:   "USD",
			TargetCurrency: "ZZZ",
		})
		require.Error(t, err, "ожидалась ошибка для несуществующей валюты")
		require.Nil(t, resp, "ответ должен быть nil для несуществующей валюты")
	})

	t.Run("empty currency codes", func(t *testing.T) {
		resp, err := s.Client.GetExchangeRate(ctx, &currency_v1.ExchangeRateRequest{
			BaseCurrency:   "",
			TargetCurrency: "",
		})
		require.Error(t, err, "ожидалась ошибка для пустых кодов валют")
		require.Nil(t, resp, "ответ должен быть nil для пустых кодов валют")
	})

	t.Run("negative amount", func(t *testing.T) {
		resp, err := s.Client.Exchange(ctx, &currency_v1.ExchangeRequest{
			BaseCurrency:   "USD",
			TargetCurrency: "EUR",
			Amount:         -100.0,
		})
		require.Error(t, err, "ожидалась ошибка для отрицательного amount")
		require.Nil(t, resp, "ответ должен быть nil для отрицательного amount")
	})
}
