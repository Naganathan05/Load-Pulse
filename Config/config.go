package Config

import (
	"log"
	"os"
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
	err := godotenv.Load("../.env");
	if err != nil {
		return err;
	}
	return nil;
}

func StringToInt(operand string) (int, error) {
	convertedInt, err := strconv.Atoi(operand);
	if err != nil {
		return -1, err;
	}
	return convertedInt, nil;
}

func GetConfig() *Config {
	configInitialized.Do(func () {
		LoadEnv();

		var clusterSize, requestSleepTime int;
		var err error;
		clusterSizeStr := os.Getenv("CLUSTER_SIZE");
		clusterSize, err = StringToInt(clusterSizeStr);
		if err != nil {
			log.Fatalf("[ERROR]: Invalid Cluster Size !!");
		}

		requestSleepTimeStr := os.Getenv("REQUEST_SLEEP_TIME");
		requestSleepTime, err = StringToInt(requestSleepTimeStr);
		if err != nil {
			log.Fatalf("[ERROR]: Invalid Request Sleep Time !!");
		}

		configInstance = Config{
			RedisKey: os.Getenv("REDIS_KEY"),
			ClusterSize: clusterSize,
			BaseQueueName: os.Getenv("BASE_QUEUE_NAME"),
			RequestSleepTime: requestSleepTime,
		}
	});
	return &configInstance;
}