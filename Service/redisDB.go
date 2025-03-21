package Service

import (
	"context"
	"fmt"
	"log"

	Config "github.com/Naganathan05/Load-Pulse/config"
	redis "github.com/redis/go-redis/v9"
)

var client *redis.Client;

var ctx = context.Background();

func IncrementRequestCount() {
	cfg := Config.GetConfig();
	err := client.Incr(ctx, cfg.RedisKey).Err();
	if err != nil {
		log.Fatal("[ERR]: Error in Incrementing Concurrent Request Count from Redis !!", err);
	}
}

func DecrementRequestCount() {
	cfg := Config.GetConfig();
	err := client.Decr(ctx, cfg.RedisKey).Err();
	if err != nil {
		log.Fatal("[ERR]: Error in Decrementing Concurrent Request Count from Redis !!", err);
	}
}

func GetRequestCount() int64 {
	cfg := Config.GetConfig();
	currentCount, err := client.Get(ctx, cfg.RedisKey).Int64();
	if err != nil {
		log.Fatal("[ERR]: Error in Fetching Concurrent Requests Count from Redis !!", err);
		return 0;
	}
	return currentCount;
}

func InitRedisClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	fmt.Println("[LOG]: Redis Client Initialized...");
}