package Raft

import (
	// "fmt"
	"time"

	"github.com/valyala/fasthttp"

	"Load-Pulse/Config"
	"Load-Pulse/Service"
	"Load-Pulse/Statistics"
)

func RunTest(workerID int, l *Service.LoadTester) *Statistics.Stats {
	start := time.Now()

	cfg := Config.GetConfig()
	requestSleepTime := cfg.RequestSleepTime

	currConcurrencyCount := Service.GetRequestCount()
	for currConcurrencyCount > int64(l.ConcurrencyLimit) {
		// workerMsg := fmt.Sprintf("[WORKER-ALERT-%d]: Concurrency Count: %d => Limit Reached !! Waiting\n", workerID, currConcurrencyCount);
		// Service.LogError(workerMsg);
		time.Sleep(time.Millisecond * time.Duration(requestSleepTime))
		currConcurrencyCount = Service.GetRequestCount()
	}

	Service.IncrementRequestCount()

	req := fasthttp.AcquireRequest()
	l.Request.CopyTo(req)
	if len(l.RequestBody) > 0 {
		req.SetBody(l.RequestBody)
	}
	resp := fasthttp.AcquireResponse()

	err := l.Client.Do(req, resp)
	rd := time.Since(start)

	stats := &Statistics.Stats{
		Endpoint:       l.Endpoint,
		ResponseDur:    rd,
		TotalRequests:  1,
		FailedRequests: 0,
	}

	if err != nil {
		stats.FailedRequests = 1
		Service.DecrementRequestCount()
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return stats
	}

	stats.ResponseSize = float64(len(resp.Body()))

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	Service.DecrementRequestCount()
	return stats
}
