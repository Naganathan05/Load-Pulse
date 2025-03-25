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
			fmt.Printf("[ERR]: Failed to Consume From Queue: %v\n", err);
			time.Sleep(1 * time.Second);
			continue;
		}

		for msg := range msgs {
			var stats Statistics.Stats;
			err := json.Unmarshal(msg.Body, &stats);
			if err != nil {
				fmt.Printf("[ERR]: Failed to unmarshal stats: %v\n", err);
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
				fmt.Printf("[AGGREGATOR]: Consumed All (%d) Events for Queue %s\n", consumedCount, queueName);
				break;
			}
		}

		fmt.Printf("[AGGREGATOR]: Consumed %d/%d Events for Queue %s\n", consumedCount, expectedEvents, queueName);

		time.Sleep(50 * time.Millisecond);
	}

	fmt.Printf("[AGGREGATOR]: Final Aggregated Stats for %s:\n", queueName);
	aggregatedStats.Print();
}