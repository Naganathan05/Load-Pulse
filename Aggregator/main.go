package main

import (
	"fmt"
	"os"
	"log"
	"sync"

	"Load-Pulse/Config"
	"Load-Pulse/Service"
	"Load-Pulse/Aggregator/Util"
)

func main() {
	Service.ConnectRabbitMQ();

	arg := os.Args[1];

	testObj, err := Service.NewLoadTester(arg);
	if err != nil {
		log.Fatal("[ERR]: Invalid File Argument:", err);
	}

	if err != nil {
		fmt.Printf("[ERROR]: Failed to load config: %v\n", err);
		os.Exit(1);
	}

	fmt.Println("[AGGREGATOR]: Starting Aggregator Service")

	var wg sync.WaitGroup;
	cfg := Config.GetConfig();

	var eventCount sync.Map;

	for testerIndex, tester := range testObj.Testers {
		queueName := fmt.Sprintf("%s-%d", cfg.BaseQueueName, testerIndex+1)

		totalRequests := tester.Conns * int(tester.Dur.Seconds()) / int(tester.Rate.Milliseconds());
		numWorkersPerCluster := Service.Min(cfg.ClusterSize, totalRequests);
		numClusters := totalRequests / numWorkersPerCluster;

		eventCount.Store(queueName, numClusters);
		fmt.Printf("[AGGREGATOR]: Listening on Queue %s With %d Expected Events\n", queueName, numClusters);

		wg.Add(1);
		go func(qName string) {
			defer wg.Done();
			Util.AggregateStatsWithCount(qName, &eventCount, tester.Endpoint);
		}(queueName);
	}

	wg.Wait();
	Service.CloseRabbitMQ();
	fmt.Println("[AGGREGATOR]: Aggregation Completed For All Endpoints.");
}