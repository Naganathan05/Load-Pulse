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

	fmt.Printf("[WORKER-%d]: Starting Worker for %s | Max Requests: %d\n", id, tester.Endpoint, maxRequests);
	ticker := time.NewTicker(tester.Rate);
	defer ticker.Stop();

	stop := time.After(tester.Dur);
	requestsMade := 0;
	stats := &Statistics.Stats{};

	for {
		select {
		case <- stop:
			fmt.Printf("[WORKER-%d]: Stopping Worker\n", id);
			leaderCh <- stats;
			return;

		case <- ticker.C:
			mu.Lock();
			if requestsMade >= maxRequests {
				mu.Unlock();
				fmt.Printf("[WORKER-%d]: Max Requests Reached. Stopping...\n", id);
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
