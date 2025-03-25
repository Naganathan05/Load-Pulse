package Util

import (
	"fmt"
	"encoding/json"
	"sync"
	"time"

	"Load-Pulse/Service"
	"Load-Pulse/Statistics"
)

func AggregateStatsWithCount(queueName string, eventCount *sync.Map, endpoint string) {

	aggregatedStats := &Statistics.Stats{Endpoint: endpoint};
	consumedCount := 0;

	expectedCount, _ := eventCount.Load(queueName);
	expectedEvents := expectedCount.(int);

	for consumedCount < expectedEvents {
		msgs, err := Service.ConsumeFromQueue(queueName);
		if err != nil {
			errMsg := fmt.Sprintf("[ERR]: Failed to Consume From Queue: %v\n", err);
			Service.LogError(errMsg);
			time.Sleep(1 * time.Second);
			continue;
		}

		for msg := range msgs {
			var stats Statistics.Stats;
			err := json.Unmarshal(msg.Body, &stats);
			if err != nil {
				errMsg := fmt.Sprintf("[ERR]: Failed to unmarshal stats: %v\n", err);
				Service.LogError(errMsg);
				continue;
			}

			aggregatedStats.Lock();
			aggregatedStats.TotalRequests += stats.TotalRequests;
			aggregatedStats.FailedRequests += stats.FailedRequests;
			aggregatedStats.ResponseSize += stats.ResponseSize;
			aggregatedStats.ResponseDur += stats.ResponseDur;
			aggregatedStats.Unlock();

			consumedCount += 1;
			eventCount.Store(queueName, consumedCount);

			if consumedCount >= expectedEvents {
				logMsg := fmt.Sprintf("[AGGREGATOR]: Consumed All (%d) Events for Queue %s\n", consumedCount, queueName);
				Service.LogCluster(logMsg);
				break;
			}
		}

		logMsg := fmt.Sprintf("[AGGREGATOR]: Consumed %d/%d Events for Queue %s\n", consumedCount, expectedEvents, queueName);
		Service.LogCluster(logMsg);

		time.Sleep(50 * time.Millisecond);
	}

	logMsg := fmt.Sprintf("[AGGREGATOR]: Final Aggregated Stats for %s:\n", queueName);
	Service.LogCluster(logMsg);
	aggregatedStats.Print();
}