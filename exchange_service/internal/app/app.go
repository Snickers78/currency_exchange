package app

import (
	"exchange_service/internal/app/grpcapp"
	exchangegrpc "exchange_service/internal/grpc/currency"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, port int, exchanger exchangegrpc.Exchanger) *App {
	grpcApp := grpcapp.NewApp(log, exchanger, port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
