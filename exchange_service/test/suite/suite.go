package suite

import (
	"context"
	currency_v1 "exchange_service/gen/proto"
	"exchange_service/internal/config"
	"net"
	"os"
	"strconv"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	Test   *testing.T
	Cfg    *config.Config
	Client currency_v1.CurrencyServiceClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	os.Setenv("CONFIG_PATH", "d:/currency-exchange/exchange_service/config/.env")
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	conn, err := grpc.NewClient(net.JoinHostPort("localhost", strconv.Itoa(cfg.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc Client error: %v", err)
	}

	return ctx, &Suite{
		Test:   t,
		Cfg:    cfg,
		Client: currency_v1.NewCurrencyServiceClient(conn),
	}
}
