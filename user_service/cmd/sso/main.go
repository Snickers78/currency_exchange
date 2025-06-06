package main

import (
	"currency-exchange/user_service/internal/app"
	"currency-exchange/user_service/internal/config"
	"currency-exchange/user_service/internal/lib/logger"
	"currency-exchange/user_service/internal/services/auth"
	storage "currency-exchange/user_service/internal/storage/postgres"
	"os"
	"os/signal"
	"syscall"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.InitLogger(EnvLocal)
	storage := storage.NewStorage(cfg.StoragePath)
	authService := auth.NewAuthService(logger, storage, cfg.TockenTTL, cfg)
	application := app.New(logger, cfg.Port, cfg.StoragePath, cfg.TockenTTL, authService)
	go application.GRPCSrv.Run()

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCSrv.Stop()
	logger.Info("Application stopped")
}
