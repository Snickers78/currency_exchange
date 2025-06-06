package suite

import (
	"context"
	ssov1 "currency-exchange/user_service/gen"
	"currency-exchange/user_service/internal/config"
	"net"
	"strconv"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	Test   *testing.T
	Cfg    *config.Config
	Client ssov1.AuthClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoad()

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
		Client: ssov1.NewAuthClient(conn),
	}
}
