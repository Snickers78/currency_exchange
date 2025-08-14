package exchange

import (
	"api_gateway/internal/config"
	exchange "api_gateway/internal/gen/exchange/proto"
	"context"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExchangeClient struct {
	client exchange.CurrencyServiceClient
}

func NewExchangeClient(cfg *config.Config) *ExchangeClient {
	// systemRoots, err := x509.SystemCertPool()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tlsCreds := credentials.NewTLS(&tls.Config{
	// 	RootCAs: systemRoots,
	// })

	conn, err := grpc.NewClient(net.JoinHostPort("localhost", strconv.Itoa(cfg.ExchangeServicePort)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &ExchangeClient{
		client: exchange.NewCurrencyServiceClient(conn),
	}
}


func (e *ExchangeClient) GetExchangeRate(ctx context.Context, req *exchange.ExchangeRateRequest) (*exchange.ExchangeRateResponse, error) {
	resp, err := e.client.GetExchangeRate(ctx, req)
	if err != nil {
		log.Println(err)
	}
	return resp, nil
}

func (e *ExchangeClient) Exchange(ctx context.Context, req *exchange.ExchangeRequest) (*exchange.ExchangeResponse, error) {
	resp, err := e.client.Exchange(ctx, req)
	if err != nil {
		log.Println(err)
	}
	return resp, nil
}
