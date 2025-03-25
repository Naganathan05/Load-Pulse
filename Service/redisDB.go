package Service

import (
	"fmt"
	"log"
	"context"

	config "Load-Pulse/Config"
	redis "github.com/redis/go-redis/v9"
)

var client *redis.Client;

var ctx = context.Background();

func IncrementRequestCount() {
	cfg := config.GetConfig();
	err := client.Incr(ctx, cfg.RedisKey).Err();
	if err != nil {
		log.Fatal("[ERR]: Error in Incrementing Concurrent Request Count from Redis !!", err);
	}
}

func DecrementRequestCount() {
	cfg := config.GetConfig();
	err := client.Decr(ctx, cfg.RedisKey).Err();
	if err != nil {
		log.Fatal("[ERR]: Error in Decrementing Concurrent Request Count from Redis !!", err);
	}
}

func GetRequestCount() int64 {
	cfg := config.GetConfig();
	currentCount, err := client.Get(ctx, cfg.RedisKey).Int64();
	if err != nil {
		log.Fatal("[ERR]: Error in Fetching Concurrent Requests Count from Redis !!", err);
		return 0;
	}
	return currentCount;
}

func ResetRequestCount() {
	cfg := config.GetConfig();
	err := client.Set(ctx, cfg.RedisKey, 0, 0).Err();
	if err != nil {
		log.Fatal("[ERR]: Error in Resetting Concurrent Requests Count in Redis !!", err);
	}
	fmt.Println("[LOG]: Concurrency Count Reset Done");
}

func InitRedisClient() {
	cfg := config.GetConfig();
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
		Protocol: 2,
	})
	fmt.Println("[LOG]: Redis Client Initialized");
}