package Tester

import (
	"fmt"
	"sync"
	"time"
	"Load-Pulse/Statistics"
)

func startWorker(id int, tester *LoadTester, leaderCh chan *Statistics.Stats, wg *sync.WaitGroup) {
	defer wg.Done();

	fmt.Printf("[WORKER-%d]: Starting worker for endpoint: %s\n", id, tester.Endpoint);

	ticker := time.NewTicker(tester.Rate);
	defer ticker.Stop();
	stop := time.After(tester.Dur);

	for {
		select {
		case <- stop:
			fmt.Printf("[WORKER-%d]: Stopping Worker\n", id);
			return;
		case <- ticker.C:
			stats := tester.RunTest();
			leaderCh <- stats;
		}
	}
}