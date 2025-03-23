package Config

import (
	"os"
	"log"
	"strconv"
	"sync"
	"github.com/joho/godotenv"
)

var (
	configInstance Config
	configInitialized sync.Once
)

type Config struct {
	RedisKey string
	ClusterSize int
	BaseQueueName string
	RequestSleepTime int
}

func LoadEnv() error {
	err := godotenv.Load(".env");
	if err != nil {
		return err;
	}
	return nil;
}

func StringToInt(operand string) int {
	convertedInt, err := strconv.Atoi(operand);
	if err != nil {
		log.Fatal("[ERR]: Invalid Integer Value In `.env` File");
		return -1;
	}
	return convertedInt;
}

func GetConfig() *Config {
	configInitialized.Do(func () {
		LoadEnv();

		clusterSizeStr := os.Getenv("CLUSTER_SIZE");
		clusterSize := StringToInt(clusterSizeStr);

		requestSleepTimeStr := os.Getenv("REQUEST_SLEEP_TIME");
		requestSleepTime := StringToInt(requestSleepTimeStr);

		configInstance = Config{
			RedisKey: os.Getenv("REDIS_KEY"),
			ClusterSize: clusterSize,
			BaseQueueName: os.Getenv("BASE_QUEUE_NAME"),
			RequestSleepTime: requestSleepTime,
		}
	});
	return &configInstance;
}