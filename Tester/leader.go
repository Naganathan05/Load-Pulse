package Tester

import (
	"fmt"
	"sync"
	"Load-Pulse/Statistics"
)

type Leader struct {
	id        int
	stats     *Statistics.Stats
	workerCnt int
}

func StartLeader(id int, tester *LoadTester, workerCnt int, wg *sync.WaitGroup, globalChan chan<- *Statistics.Stats) {
	defer wg.Done();

	leader := &Leader{
		id:        id,
		stats:     &Statistics.Stats{Endpoint: tester.Endpoint},
		workerCnt: workerCnt,
	}

	leaderChan := make(chan *Statistics.Stats, workerCnt);
	var workerWg sync.WaitGroup;

	for i := range leader.workerCnt {
		workerWg.Add(1);
		go startWorker(i, tester, leaderChan, &workerWg);
	}

	go func() {
		workerWg.Wait();
		close(leaderChan);
	}();

	for stats := range leaderChan {
		leader.stats.Lock();
		leader.stats.TotalRequests += stats.TotalRequests;
		leader.stats.FailedRequests += stats.FailedRequests;
		leader.stats.ResponseSize += stats.ResponseSize;
		leader.stats.ResponseDur += stats.ResponseDur;
		leader.stats.Unlock();
	}

	fmt.Printf("[LEADER-%d]: Sending final stats to global aggregator\n", leader.id);
	globalChan <- leader.stats;
}