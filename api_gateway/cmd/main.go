package main

import (
	"api_gateway/internal/auth"
	"api_gateway/internal/config"
	"api_gateway/internal/exchange"
	"api_gateway/internal/kafka"
	"api_gateway/internal/middleware"
	authHandler "api_gateway/internal/server/auth_handler"
	exchangeHandler "api_gateway/internal/server/exchange_handler"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var (
	brokers = []string{"localhost:9093"}
	topic   = "logs"
)

func main() {
	router := gin.Default()
	cfg := config.LoadConfig()
	// log := logger.InitLogger(EnvLocal)
	kafkaHook := kafka.NewKafkaHook(brokers, topic)

	//clients
	authClient := auth.NewAuthCLient(cfg)
	exchangeClient := exchange.NewExchangeClient(cfg)

	//handlers
	_ = authHandler.NewAuthHandler(router, authClient, kafkaHook)
	_ = exchangeHandler.NewExchangeHandler(router, exchangeClient, kafkaHook)

	router.Use(middleware.CORS())

	if err := router.Run(":" + strconv.Itoa(cfg.GatewayPort)); err != nil {
		log.Fatal("Failed to start server", "err", err)
	}
}
