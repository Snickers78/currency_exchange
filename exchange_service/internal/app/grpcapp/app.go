package grpcapp

import (
	exchangegrpc "exchange_service/internal/grpc/currency"
	"fmt"
	"log/slog"
	"net"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, exchanger exchangegrpc.Exchanger, port int) *App {
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.EnableClientHandlingTimeHistogram()

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)

	exchangegrpc.Register(gRPCServer, exchanger)

	return &App{log: log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic(err)
	}

	a.log.Info("GRPC server listening", slog.String("port", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		panic(err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("Stopping GRPC server", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
