package Load_Tester

import (
	"fmt"
	"sync"
	"encoding/json"
	"loadpulse.local/Statistics"
	"loadpulse.local/Service"
)

type Leader struct {
	id        int
	stats     *Statistics.Stats
	workerCnt int
	queueName string
}

func StartLeader(id int, tester *LoadTester, workerCnt int, maxRequests int, queueName string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done();

	leader := &Leader{
		id:        id,
		stats:     &Statistics.Stats{Endpoint: tester.Endpoint},
		workerCnt: workerCnt,
		queueName: queueName,
	}

	leaderChan := make(chan *Statistics.Stats, workerCnt);
	var workerWg sync.WaitGroup;

	for i := 0; i < leader.workerCnt; i++ {
		workerWg.Add(1);
		go startWorker(i, tester, leaderChan, &workerWg, mu, maxRequests);
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

	statsJSON, _ := json.Marshal(leader.stats);
	fmt.Printf("[LEADER-%d]: Publishing Stats to Queue: %s\n", leader.id, queueName);
	err := Service.PublishToQueue(queueName, statsJSON);
	if err != nil {
		fmt.Printf("[ERROR]: Failed to publish stats: %v\n", err);
	}
}