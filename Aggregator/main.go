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
		log.Fatal("[ERR]: Invalid File Arguement:", err);
	}

	if err != nil {
		fmt.Printf("[ERROR]: Failed to load config: %v\n", err);
		os.Exit(1);
	}

	fmt.Println("[AGGREGATOR]: Starting Aggregator Service...");

	var wg sync.WaitGroup;
	cfg := Config.GetConfig();

	for testerIndex, tester := range testObj.Testers {
		queueName := fmt.Sprintf("%s-%d", cfg.BaseQueueName, testerIndex + 1);

		fmt.Println("Total Conns:", tester.Conns);
		fmt.Println("Duration: ", int(tester.Dur.Seconds()));
		fmt.Println("Rate: ", int(tester.Rate.Milliseconds()));

		totalRequests := tester.Conns * int(tester.Dur.Seconds()) / int(tester.Rate.Milliseconds());
		numWorkersPerCluster := Service.Min(cfg.ClusterSize, totalRequests);
		numClusters := totalRequests / numWorkersPerCluster;

		fmt.Printf("[AGGREGATOR]: Listening on Queue %s With %d Clusters\n", queueName, numClusters);

		wg.Add(1);
		go func(qName string, clusters int) {
			defer wg.Done();
			Util.AggregateStats(qName, clusters);
		}(queueName, numClusters);
	}

	wg.Wait();
	Service.CloseRabbitMQ();
	fmt.Println("[AGGREGATOR]: Aggregation Completed For All Endpoints.");
}