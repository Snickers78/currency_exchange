package grpcapp

import (
	authgrpc "currency-exchange/user_service/internal/grpc/auth"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

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
