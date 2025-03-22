package Tester

import (
	"io"
	"net/http"
	"time"
	"Load-Pulse/Service"
	"Load-Pulse/Statistics"
)

type LoadTester struct {
	Endpoint string
	Conns    int
	Request  *http.Request
	Client   *http.Client
	Stats    *Statistics.Stats
	Dur      time.Duration
	Rate     time.Duration
}

func NewTester(r *http.Request, conns int, dur, rate time.Duration, end string) *LoadTester {
	return &LoadTester{
		Endpoint: end,
		Request:  r,
		Client:   &http.Client{},
		Conns:    conns,
		Dur:      dur,
		Rate:     rate,
		Stats:    &Statistics.Stats{Endpoint: end},
	}
}

func (l *LoadTester) RunTest() *Statistics.Stats {
	var body []byte
	start := time.Now()
	Service.IncrementRequestCount()

	resp, err := l.Client.Do(l.Request)
	rd := time.Since(start)

	stats := &Statistics.Stats{
		Endpoint:      l.Endpoint,
		ResponseDur:   rd,
		TotalRequests: 1,
		FailedRequests: 0,
	}

	if err != nil {
		stats.FailedRequests = 1
		Service.DecrementRequestCount()
		return stats
	}

	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	stats.ResponseSize = float64(len(body))
	Service.DecrementRequestCount()

	return stats
}