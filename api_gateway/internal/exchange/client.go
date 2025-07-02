package exchange

import (
	"api_gateway/internal/config"
	exchange "api_gateway/internal/gen/exchange/proto"
	"context"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExchangeClient struct {
	client exchange.CurrencyServiceClient
}

func NewExchangeClient(cfg *config.Config) *ExchangeClient {
	conn, err := grpc.NewClient(net.JoinHostPort("localhost", strconv.Itoa(cfg.ExchangeServicePort)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &ExchangeClient{
		client: exchange.NewCurrencyServiceClient(conn),
	}
}

func (e *ExchangeClient) GetExchangeRate(ctx context.Context, req *exchange.ExchangeRateRequest) (*exchange.ExchangeRateResponse, error) {
	return e.client.GetExchangeRate(ctx, req)
}

func (e *ExchangeClient) Exchange(ctx context.Context, req *exchange.ExchangeRequest) (*exchange.ExchangeResponse, error) {
	return e.client.Exchange(ctx, req)
}
