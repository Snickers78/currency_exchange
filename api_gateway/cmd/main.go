package main

import (
	"api_gateway/infra/kafka"
	"api_gateway/internal/auth"
	"api_gateway/internal/config"
	"api_gateway/internal/exchange"
	"api_gateway/internal/middleware"
	authHandler "api_gateway/internal/server/auth_handler"
	exchangeHandler "api_gateway/internal/server/exchange_handler"
	"context"
	"time"

	//"context"
	"log"
	"strconv"

	//"time"

	//"time"

	"github.com/gin-gonic/gin"
)

const (
	// EnvLocal = "local"
	// EnvDev   = "dev"
	// EnvProd  = "prod"
	topic = "logs"
	App   = false
)

var (
	brokers = []string{"localhost:9093"}
)

func main() {
	//redisConnect := infraRedis.NewRedisCLient("localhost:6379")
	rateLimiter := middleware.NewBucketLimiter(context.Background(), 100, 1*time.Second)
	router := gin.Default()
	cfg := config.LoadConfig(App)
	// log := logger.InitLogger(EnvLocal)
	kafkaHook := kafka.NewKafkaHook(brokers, topic)

	//clients
	authClient := auth.NewAuthClient(cfg)
	exchangeClient := exchange.NewExchangeClient(cfg)

	//middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	//handlers
	authHandler.NewAuthHandler(router, authClient, kafkaHook, cfg.Secret)
	exchangeHandler.NewExchangeHandler(router, exchangeClient, kafkaHook, cfg.Secret)

	if err := router.Run(":" + strconv.Itoa(cfg.GatewayPort)); err != nil {
		log.Fatal("Failed to start server", "err", err)
	}
}
