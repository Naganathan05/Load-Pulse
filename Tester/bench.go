package Tester

import (
	"io"
	"fmt"
	"sync"
	"time"
	"math"
	"bytes"
	"net/http"

	"Load-Pulse/Statistics"
	"Load-Pulse/Config"
)

type Bench struct {
	testers []*LoadTester
	ch      chan *Statistics.Stats
}

func New(path string) (*Bench, error) {
	var testers []*LoadTester;

	conf, err := fromJSON(path)
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

	numWorkersPerCluster := cfg.ClusterSize;

	fmt.Println("[LOG]: Starting Load Test for Individual Endpoints By Clustering...");
	globalStatsChan := make(chan *Statistics.Stats, len(b.testers));

	for testerIndex, tester := range b.testers {

		totalRequests := tester.Conns * int(tester.Rate.Milliseconds()) * int(tester.Dur.Seconds());
		numClusters := int(math.Ceil(float64(totalRequests) / float64(numWorkersPerCluster)));

		for clusterID := 0; clusterID < numClusters; clusterID++ {
			wg.Add(1);

			go func(t *LoadTester, clusterID, testerIndex int) {
				defer wg.Done();

				fmt.Printf("[Cluster-%d, Tester-%d]: Initializing Leader and Worker Nodes...\n", clusterID + 1, testerIndex + 1);
				StartLeader(clusterID, t, numWorkersPerCluster, &wg, globalStatsChan);

			}(tester, clusterID, testerIndex);
		}
	}

	wg.Wait();
	close(globalStatsChan);
	fmt.Println("\n[GLOBAL]: Aggregating Stats From All Clusters...");

	finalStats := &Statistics.Stats{}
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