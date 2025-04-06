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
		errMsg := fmt.Sprintf("[ERROR]: Failed to load config: %v\n", err);
		Service.LogError(errMsg);
		os.Exit(1);
	}

	Service.LogCluster("[AGGREGATOR]: Starting Aggregator Service\n")

	var wg sync.WaitGroup;
	cfg := Config.GetConfig();

	var eventCount sync.Map;

	for testerIndex, tester := range testObj.Testers {
		queueName := fmt.Sprintf("%s-%d", cfg.BaseQueueName, testerIndex+1)

		totalRequests := tester.Conns * int(tester.Dur.Seconds()) / int(tester.Rate.Milliseconds());
		numWorkersPerCluster := Service.Min(cfg.ClusterSize, totalRequests);
		numClusters := totalRequests / numWorkersPerCluster;

		eventCount.Store(queueName, numClusters);
		logMsg := fmt.Sprintf("[AGGREGATOR]: Starting Aggregator Service for %s with %d clusters\n", queueName, numClusters);
		Service.LogCluster(logMsg);

		wg.Add(1);
		go func(qName string) {
			defer wg.Done();
			Util.AggregateStatsWithCount(qName, &eventCount, tester.Endpoint);
		}(queueName);
	}

	wg.Wait();
	Service.CloseRabbitMQ();
	Service.LogCluster("[AGGREGATOR]: Aggregation Completed For All Endpoints.\n");
}