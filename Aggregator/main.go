package Aggregator

import (
	"fmt"
	"os"
	"sync"

	"loadpulse.local/Config"
	"loadpulse.local/Service"
	"loadpulse.local/Load_Tester"
)

func min(a int, b int) int {
	if a < b {
		return a;
	}
	return b;
}

func main() {
	Service.ConnectRabbitMQ();
	defer Service.CloseRabbitMQ();

	arg := os.Args[1];

	config, err := Load_Tester.FromJSON(arg);
	if err != nil {
		fmt.Printf("[ERROR]: Failed to load config: %v\n", err);
		os.Exit(1);
	}

	fmt.Println("[AGGREGATOR]: Starting Aggregator Service...");

	var wg sync.WaitGroup;
	cfg := Config.GetConfig();

	for testerIndex, req := range config.Req {
		queueName := fmt.Sprintf("%s-%d", cfg.BaseQueueName, testerIndex + 1);

		totalRequests := req.Connections * int(config.Duration.Seconds()) / int(req.Rate.Seconds());
		numWorkersPerCluster := min(cfg.ClusterSize, totalRequests);
		numClusters := totalRequests / numWorkersPerCluster;

		fmt.Printf("[AGGREGATOR]: Listening on Queue %s With %d Clusters\n", queueName, numClusters);
		
		wg.Add(1);
		go func(qName string, clusters int) {
			defer wg.Done();
			AggregateStats(qName, clusters);
		}(queueName, numClusters);
	}

	wg.Wait();
	fmt.Println("[AGGREGATOR]: Aggregation Completed For All Endpoints.");
}