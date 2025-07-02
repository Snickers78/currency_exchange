package main

import (
	"os"
	"os/signal"
	"syscall"
	"user_service/internal/app"
	"user_service/internal/config"
	"user_service/internal/lib/logger"
	"user_service/internal/metrics"
	"user_service/internal/services/auth"
	storage "user_service/internal/storage/postgres"
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
	metricsApp := metrics.NewMetricsApp()
	go application.GRPCSrv.Run()
	go metricsApp.Run()

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCSrv.Stop()
	metricsApp.Stop()
	logger.Info("Application stopped")
}
