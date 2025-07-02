package main

import (
	"exchange_service/internal/app"
	"exchange_service/internal/config"
	"exchange_service/internal/lib/logger"
	"exchange_service/internal/metrics"
	currency "exchange_service/internal/services"
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
	cfg := config.LoadConfig()
	logger := logger.InitLogger(EnvLocal)
	exchanger := currency.NewExchanger(logger, cfg)
	application := app.New(logger, cfg.Port, exchanger)
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
