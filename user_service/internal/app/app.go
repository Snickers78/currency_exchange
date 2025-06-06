package app

import (
	"currency-exchange/user_service/internal/app/grpcapp"
	authgrpc "currency-exchange/user_service/internal/grpc/auth"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath string, tockenTTL time.Duration, authService authgrpc.Auth) *App {
	grpcApp := grpcapp.NewApp(log, authService, port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
