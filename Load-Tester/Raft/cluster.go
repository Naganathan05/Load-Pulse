package Raft

import (
	"log"
	"fmt"
	"sync"

	"Load-Pulse/Config"
	"Load-Pulse/Service"
)

func Run(b *Service.Bench) {
	var wg sync.WaitGroup;

	cfg := Config.GetConfig();

	fmt.Println("[LOG]: Starting Load Test for Individual Endpoints By Clustering");

	var mu sync.Mutex;
	for testerIndex, tester := range b.Testers {
		totalRequests := tester.Conns * int(tester.Dur.Seconds()) / int(tester.Rate.Milliseconds());
		numWorkersPerCluster := Service.Min(cfg.ClusterSize, totalRequests);
		numClusters := totalRequests / numWorkersPerCluster;

		requestsPerWorker := totalRequests / numWorkersPerCluster;
		remainingRequests := totalRequests % numWorkersPerCluster;

		baseQueueName := cfg.BaseQueueName;
		queueName := fmt.Sprintf("%s-%d", baseQueueName, testerIndex + 1);

		/* ------------------------   DEBUGGING  --------------------------------
			fmt.Println("Total Requests:", totalRequests);
		 	fmt.Println("Number of Clusters:", numClusters);
		 	fmt.Println("Number of Workers Per Cluster:", numWorkersPerCluster);
			fmt.Println("Number of Requests Per Worker:", requestsPerWorker);
		 	fmt.Println("Number of Remaining Requests:", remainingRequests);
		 	fmt.Println("Number of Connections Required:", tester.Conns);
		 	fmt.Println("Concurrency Limit:", tester.ConcurrencyLimit);
		 	fmt.Println("Request Rate:", int(tester.Rate.Milliseconds()));
		---------------------------------------------------------------------------*/

		err := Service.DeleteQueue(queueName);
		if err != nil {	
			log.Fatalf("[ERROR]: Failed to Delete Queue: %v\n", err);
		}

		fmt.Printf("[LOG]: Tester %d → Total Requests: %d | Workers: %d | Req/Worker: %d | Remaining: %d\n",
			testerIndex+1, totalRequests, numWorkersPerCluster, requestsPerWorker, remainingRequests);

		for clusterID := 0; clusterID < numClusters; clusterID++ {
			wg.Add(1);

			finalRequests := requestsPerWorker;
			if clusterID == numWorkersPerCluster - 1 {
				finalRequests += remainingRequests;
			}

			go func(t *Service.LoadTester, clusterID, testerIndex, finalRequests int) {

				fmt.Printf("[Cluster-%d, Tester-%d]: Starting Leader with %d Requests\n", clusterID + 1, testerIndex + 1, finalRequests);
				StartLeader(clusterID, t, numWorkersPerCluster, finalRequests, queueName, &wg, &mu);

			}(tester, clusterID, testerIndex, finalRequests);
		}
	}
	wg.Wait();
}