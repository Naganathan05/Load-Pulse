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
	RedisURL string
	RedisPassword string
	RabbitMQURL string
	TesterServerPort string
	AggregatorServerPort string
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
			RedisURL: os.Getenv("REDIS_URL"),
			RabbitMQURL: os.Getenv("RABBITMQ_URL"),
			ClusterSize: clusterSize,
			RedisKey: os.Getenv("REDIS_KEY"),
			RequestSleepTime: requestSleepTime,
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
			BaseQueueName: os.Getenv("BASE_QUEUE_NAME"),
			TesterServerPort: os.Getenv("LOAD_TESTER_PORT"),
			AggregatorServerPort: os.Getenv("AGGREGATOR_PORT"),
		}
	});
	return &configInstance;
}