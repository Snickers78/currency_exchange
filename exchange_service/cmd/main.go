package main

import (
	"exchange_service/infra/metrics"
	"exchange_service/internal/app"
	"exchange_service/internal/config"
	logg "exchange_service/internal/lib/logger"
	currency "exchange_service/internal/services"
	"os"
	"os/signal"
	"syscall"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
	//topic    = "logs"
)

func main() {
	cfg := config.LoadConfig()
	//brokers := []string{cfg.Broker1, cfg.Broker2}
	//kafkaHook := kafka.NewKafkaHook(brokers, topic)
	logger := logg.InitLogger(EnvLocal)
	exchanger := currency.NewExchanger(cfg, logger)
	application := app.New(logger, cfg.Port, exchanger)
	metricsApp := metrics.NewMetricsApp()
	go application.GRPCSrv.Run()
	go metricsApp.Run()

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	metricsApp.Stop()
	application.GRPCSrv.Stop()
	logger.Info("Application stopped")
}
