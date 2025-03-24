package main

import (
	"fmt"
	"encoding/json"
	
	"Load-Pulse/Service"
	"Load-Pulse/Statistics"
)

func AggregateStats(queueName string, expectedClusters int) {
	fmt.Printf("[AGGREGATOR]: Listening on Queue: %s\n", queueName)

	msgs, err := Service.ConsumeFromQueue(queueName)
	if err != nil {
		fmt.Printf("[ERR]: Failed to Consume From Queue: %v\n", err)
		return
	}

	aggregatedStats := &Statistics.Stats{}
	clusterCount := 0

	for msg := range msgs {
		var stats Statistics.Stats
		err := json.Unmarshal(msg.Body, &stats)
		if err != nil {
			fmt.Printf("[ERR]: Failed to unmarshal stats: %v\n", err)
			continue
		}

		aggregatedStats.Lock()
		aggregatedStats.TotalRequests += stats.TotalRequests
		aggregatedStats.FailedRequests += stats.FailedRequests
		aggregatedStats.ResponseSize += stats.ResponseSize
		aggregatedStats.ResponseDur += stats.ResponseDur
		aggregatedStats.Unlock()

		clusterCount += 1
		if clusterCount >= expectedClusters {
			break
		}
	}

	fmt.Printf("\n[AGGREGATOR]: Final Aggregated Stats for %s:\n", queueName)
	aggregatedStats.Print()
}
