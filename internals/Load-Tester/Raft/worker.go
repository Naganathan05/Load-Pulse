package Raft

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"Load-Pulse/Service"
	"Load-Pulse/Statistics"
)

func startWorker(id int, tester *Service.LoadTester, leaderCh chan *Statistics.Stats, wg *sync.WaitGroup, mu *sync.Mutex, maxRequests int) {
	defer wg.Done();

	workerMsg := fmt.Sprintf("[WORKER-%d]: Starting Worker for %s | Max Requests: %d\n", id, tester.Endpoint, maxRequests);
	Service.LogWorker(workerMsg);
	ticker := time.NewTicker(tester.Rate);
	defer ticker.Stop();

	stop := time.After(tester.Dur);
	requestsMade := 0;
	stats := &Statistics.Stats{MinResponseTime: time.Second * 1000};

	for {
		select {
		case <- stop:
			// workerMsg := fmt.Sprintf("[WORKER-%d]: Stopping Worker\n", id);
			// Service.LogWorker(workerMsg);
			leaderCh <- stats;
			return;

		case <- ticker.C:
			mu.Lock();
			if requestsMade >= maxRequests {
				mu.Unlock();
				// workerMsg := fmt.Sprintf("[WORKER-%d]: Load Requests Done. Terminating\n", id);
				// Service.LogWorker(workerMsg);
				leaderCh <- stats;
				return;
			}

			requestsMade += 1;
			mu.Unlock();

			newStats := RunTest(id, tester);
			if newStats.FailedRequests > 0 {
				stats.Update(int(newStats.ResponseSize), newStats.ResponseDur, errors.New("request failed"));
			} else {
				stats.Update(int(newStats.ResponseSize), newStats.ResponseDur, nil);
			}
		}
	}
}
