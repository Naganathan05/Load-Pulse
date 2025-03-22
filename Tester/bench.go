package Tester

import (
	"Load-Pulse/Statistics"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Bench controls the testers, cluster creation, and stat aggregation
type Bench struct {
	testers []*LoadTester
	ch      chan *Statistics.Stats
}

// New returns a Bench tester
func New(path string) (*Bench, error) {
	var testers []*LoadTester

	conf, err := fromJSON(path)
	if err != nil {
		return nil, err
	}

	for _, req := range conf.Req {
		var buf io.Reader
		addr := conf.Host + req.Endpoint

		if req.Data != "" {
			buf = bytes.NewBufferString(req.Data)
		}

		r, err := http.NewRequest(req.Method, addr, buf)
		if err != nil {
			return nil, err
		}

		lt := NewTester(r, req.Connections, conf.Duration*time.Second, req.Rate*time.Millisecond, req.Endpoint)
		testers = append(testers, lt)
	}

	b := &Bench{
		testers: testers,
		ch:      make(chan *Statistics.Stats, len(testers)),
	}

	return b, nil
}

// Run performs the benchmark test with clusters
func (b *Bench) Run() {
	var wg sync.WaitGroup

	// Define the number of clusters and workers per cluster
	numClusters := 3          // Fixed number of clusters
	numWorkersPerCluster := 5 // Fixed number of worker nodes in each cluster

	// Channel for collecting final stats from all clusters
	globalStatsChan := make(chan *Statistics.Stats, len(b.testers)*numClusters)

	fmt.Println("Starting load test with clusters...")

	// Iterate over each tester and create clusters
	for testerIndex, tester := range b.testers {
		for clusterID := 0; clusterID < numClusters; clusterID++ {
			wg.Add(1)

			go func(t *LoadTester, clusterID, testerIndex int) {
				defer wg.Done()

				fmt.Printf("[Cluster-%d, Tester-%d]: Initializing leader and workers...\n", clusterID+1, testerIndex+1)

				// Start the leader node with workers in a cluster
				StartLeader(clusterID, t, numWorkersPerCluster, &wg, globalStatsChan)

			}(tester, clusterID, testerIndex)
		}
	}

	// Wait for all clusters to finish
	wg.Wait()
	close(globalStatsChan)

	// Aggregate and print final stats
	fmt.Println("\n[GLOBAL]: Aggregating stats from all clusters...");

	finalStats := &Statistics.Stats{}
	for clusterStats := range globalStatsChan {
		finalStats.Lock()
		finalStats.TotalRequests += clusterStats.TotalRequests
		finalStats.FailedRequests += clusterStats.FailedRequests
		finalStats.ResponseSize += clusterStats.ResponseSize
		finalStats.ResponseDur += clusterStats.ResponseDur
		finalStats.Unlock()
	}

	fmt.Println("\n[GLOBAL]: Final aggregated stats across all clusters:")
	finalStats.Print()
}