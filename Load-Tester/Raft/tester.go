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
	var body []byte;
	start := time.Now();

	cfg := Config.GetConfig();
	requestSleepTime := cfg.RequestSleepTime;

	currConcurrencyCount := Service.GetRequestCount();
	for currConcurrencyCount > int64(l.ConcurrencyLimit) {
		// workerMsg := fmt.Sprintf("[WORKER-ALERT-%d]: Concurrency Count: %d => Limit Reached !! Waiting\n", workerID, currConcurrencyCount);
		// Service.LogError(workerMsg);
		time.Sleep(time.Millisecond * time.Duration(requestSleepTime));
		currConcurrencyCount = Service.GetRequestCount();
	}

	Service.IncrementRequestCount();

	resp, err := l.Client.Do(l.Request);
	rd := time.Since(start);

	stats := &Statistics.Stats{
		Endpoint:       l.Endpoint,
		ResponseDur:    rd,
		TotalRequests:  1,
		FailedRequests: 0,
	}

	if err != nil {
		stats.FailedRequests = 1;
		Service.DecrementRequestCount();
		return stats;
	}

	defer resp.Body.Close();

	body, _ = io.ReadAll(resp.Body);
	stats.ResponseSize = float64(len(body));
	Service.DecrementRequestCount();
	return stats;
}