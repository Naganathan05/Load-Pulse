package Config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	configInstance Config
	configInitialized sync.Once
)

type Config struct {
	RedisKey string
	NodeID string
	RaftAddress string
}

func LoadEnv() error {
	err := godotenv.Load(".env");
	if err != nil {
		return err;
	}
	return nil;
}

func GetConfig() *Config {
	configInitialized.Do(func () {
		LoadEnv();

		configInstance = Config{
			RedisKey: os.Getenv("REDIS_KEY"),
			NodeID: os.Getenv("NODE_ID"),
			RaftAddress: os.Getenv("RAFT_ADDRESS"),
		}
	});
	return &configInstance;
}