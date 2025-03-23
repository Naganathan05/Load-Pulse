package Tester

import (
	"io"
	"fmt"
	"sync"
	"time"
	"bytes"
	"net/http"

	"Load-Pulse/Statistics"
	"Load-Pulse/Config"
)

type Bench struct {
	testers []*LoadTester
	ch      chan *Statistics.Stats
}

func min(a int, b int) int {
	if a < b {
		return a;
	}
	return b;
}

func New(path string) (*Bench, error) {
	var testers []*LoadTester;

	conf, err := fromJSON(path);
	if err != nil {
		return nil, err;
	}

	for _, req := range conf.Req {
		var buf io.Reader;
		addr := conf.Host + req.Endpoint;

		if req.Data != "" {
			buf = bytes.NewBufferString(req.Data);
		}

		r, err := http.NewRequest(req.Method, addr, buf);
		if err != nil {
			return nil, err;
		}

		lt := NewTester(r, req.Connections, conf.Duration*time.Second, req.Rate*time.Millisecond, req.Endpoint, req.ConcurrencyLimit);
		testers = append(testers, lt);
	}

	b := &Bench{
		testers: testers,
		ch:      make(chan *Statistics.Stats, len(testers)),
	}

	return b, nil;
}

func (b *Bench) Run() {
	var wg sync.WaitGroup;

	cfg := Config.GetConfig();

	fmt.Println("[LOG]: Starting Load Test for Individual Endpoints By Clustering...");
	globalStatsChan := make(chan *Statistics.Stats, len(b.testers));

	var mu sync.Mutex;
	for testerIndex, tester := range b.testers {
		totalRequests := tester.Conns * int(tester.Rate.Milliseconds()) * int(tester.Dur.Seconds());
		numWorkersPerCluster := min(cfg.ClusterSize, totalRequests);
		numClusters := totalRequests / numWorkersPerCluster;

		requestsPerWorker := totalRequests / numWorkersPerCluster;
		remainingRequests := totalRequests % numWorkersPerCluster;

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

		fmt.Printf("[LOG]: Tester %d → Total Requests: %d | Workers: %d | Req/Worker: %d | Remaining: %d\n",
			testerIndex + 1, totalRequests, numWorkersPerCluster, requestsPerWorker, remainingRequests);

		for clusterID := 0; clusterID < numClusters; clusterID++ {
			wg.Add(1);

			finalRequests := requestsPerWorker;
			if clusterID == numWorkersPerCluster - 1 {
				finalRequests += remainingRequests;
			}

			go func(t *LoadTester, clusterID, testerIndex, finalRequests int) {

				fmt.Printf("[Cluster-%d, Tester-%d]: Starting Leader with %d Requests...\n", clusterID + 1, testerIndex + 1, finalRequests);
				StartLeader(clusterID, t, numWorkersPerCluster, finalRequests, &wg, globalStatsChan, &mu);

			}(tester, clusterID, testerIndex, finalRequests);
		}
	}

	wg.Wait();
	close(globalStatsChan);
	fmt.Println("\n[GLOBAL]: Aggregating Stats From All Clusters...");

	finalStats := &Statistics.Stats{};
	for clusterStats := range globalStatsChan {
		finalStats.Lock();
		finalStats.TotalRequests += clusterStats.TotalRequests;
		finalStats.FailedRequests += clusterStats.FailedRequests;
		finalStats.ResponseSize += clusterStats.ResponseSize;
		finalStats.ResponseDur += clusterStats.ResponseDur;
		finalStats.Unlock();
	}

	fmt.Println("\n[GLOBAL]: Final aggregated stats across all clusters:");
	finalStats.Print();
}