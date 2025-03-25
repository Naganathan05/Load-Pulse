package Service

import (
	"fmt"
	"time"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	White  = "\033[37m"
	Reset  = "\033[0m"
	Violet    = "\033[35m"
)

func LogError(message string) {
	fmt.Printf("%s%s%s", Red, message, Reset)
}

func LogLeader(message string) {
	fmt.Printf("%s%s%s", Violet, message, Reset)
}

func LogWorker(message string) {
	fmt.Printf("%s%s%s", Green, message, Reset)
}

func LogCluster(message string) {
	fmt.Printf("%s%s%s", Blue, message, Reset)
}

func LogServer(message string) {
	fmt.Printf("%s%s%s", White, message, Reset)
}

func LogWithTimestamp(level, message, color string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s[%s] %s: %s %s", color, timestamp, level, message, Reset)
}