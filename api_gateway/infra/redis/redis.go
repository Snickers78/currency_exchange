package infraRedis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisCLient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Cannot connect to redis DB: %v", err))
	}

	return rdb
}
