package Service

import (
	"context"
	_ "embed"
	"log"

	config "Load-Pulse/Config"

	redis "github.com/redis/go-redis/v9"
)

var client *redis.Client

var ctx = context.Background()

//go:embed scripts/increment.lua
var kaimScript string

// Lua script to atomically check and increment
var incrementScript = redis.NewScript(kaimScript)

func TryIncrementRequestCount(limit int) (bool, error) {
	cfg := config.GetConfig()
	result, err := incrementScript.Run(ctx, client, []string{cfg.RedisKey}, limit).Int64()
	if err != nil {
		return false, err
	}
	// If result is -1, limit reached. Otherwise (new count), success.
	return result != -1, nil
}

func DecrementRequestCount() {
	cfg := config.GetConfig()
	err := client.Decr(ctx, cfg.RedisKey).Err()
	if err != nil {
		log.Fatal("[ERR]: Error in Decrementing Concurrent Request Count from Redis !!", err)
	}
}

func GetRequestCount() int64 {
	cfg := config.GetConfig()
	currentCount, err := client.Get(ctx, cfg.RedisKey).Int64()
	if err != nil {
		log.Fatal("[ERR]: Error in Fetching Concurrent Requests Count from Redis !!", err)
		return 0
	}
	return currentCount
}

func ResetRequestCount() {
	cfg := config.GetConfig()
	err := client.Set(ctx, cfg.RedisKey, 0, 0).Err()
	if err != nil {
		log.Fatal("[ERR]: Error in Resetting Concurrent Requests Count in Redis !!", err)
	}
	LogServer("[LOG]: Concurrency Count Reset Done\n")
}

func InitRedisClient() {
	cfg := config.GetConfig()
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
		Protocol: 2,
	})
	LogServer("[LOG]: Redis Client Initialized\n")
}
