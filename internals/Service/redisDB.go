package Service

import (
	"context"
	"log"
	"sync"

	config "Load-Pulse/Config"

	redis "github.com/redis/go-redis/v9"
)

var client *redis.Client;

var ctx = context.Background();

var mu sync.Mutex;

func IncrementRequestCount() {
	mu.Lock()
	defer mu.Unlock()
	cfg := config.GetConfig()
	err := client.Incr(ctx, cfg.RedisKey).Err()
	if err != nil {
		log.Fatal("[ERR]: Error in Incrementing Concurrent Request Count from Redis !!", err);
	}
}

func DecrementRequestCount() {
	mu.Lock()
    defer mu.Unlock()
	cfg := config.GetConfig();
	err := client.Decr(ctx, cfg.RedisKey).Err();
	if err != nil {
		log.Fatal("[ERR]: Error in Decrementing Concurrent Request Count from Redis !!", err);
	}
}

func GetRequestCount() int64 {
	mu.Lock()
    defer mu.Unlock()
	cfg := config.GetConfig();
	currentCount, err := client.Get(ctx, cfg.RedisKey).Int64();
	if err != nil {
		log.Fatal("[ERR]: Error in Fetching Concurrent Requests Count from Redis !!", err);
		return 0;
	}
	return currentCount;
}

func ResetRequestCount() {
	mu.Lock()
    defer mu.Unlock()
	cfg := config.GetConfig();
	err := client.Set(ctx, cfg.RedisKey, 0, 0).Err();
	if err != nil {
		log.Fatal("[ERR]: Error in Resetting Concurrent Requests Count in Redis !!", err);
	}
	LogServer("[LOG]: Concurrency Count Reset Done\n");
}

func InitRedisClient() {
	cfg := config.GetConfig();
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
		Protocol: 2,
	})
	LogServer("[LOG]: Redis Client Initialized\n");
}