package Raft

import (
	"io"
	// "fmt"
	"time"

	"Load-Pulse/Config"
	"Load-Pulse/Service"
	"Load-Pulse/Statistics"
)

func RunTest(workerID int, l *Service.LoadTester) *Statistics.Stats {
	var body []byte
	start := time.Now()

	cfg := Config.GetConfig()
	requestSleepTime := cfg.RequestSleepTime

	// Atomic Concurrency Check
	for {
		allowed, err := Service.TryIncrementRequestCount(l.ConcurrencyLimit)
		if err != nil {
			// If Redis fails, log and retry
			Service.LogError("[ERR]: Redis Error in Concurrency Check: " + err.Error() + "\n")
			time.Sleep(time.Millisecond * time.Duration(requestSleepTime))
			continue
		}
		if allowed {
			break
		}
		// Limit reached, wait and retry
		time.Sleep(time.Millisecond * time.Duration(requestSleepTime))
	}

	resp, err := l.Client.Do(l.Request)
	rd := time.Since(start)

	stats := &Statistics.Stats{
		Endpoint:       l.Endpoint,
		ResponseDur:    rd,
		TotalRequests:  1,
		FailedRequests: 0,
	}

	if err != nil {
		stats.FailedRequests = 1
		Service.DecrementRequestCount()
		return stats
	}

	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	stats.ResponseSize = float64(len(body))
	Service.DecrementRequestCount()
	return stats
}
