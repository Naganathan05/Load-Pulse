package Tester

import (
	"io"
	"net/http"
	"sync"
	"time"
	"log"
    redisDB "github.com/Naganathan05/Load-Pulse/Service"
)

type LoadTester struct {
	endpoint string
	conns    int
	request  *http.Request
	client   *http.Client
	stats    *Stats
	dur      time.Duration
	rate     time.Duration
}

func NewTester(r *http.Request, conns int, dur, rate time.Duration, end string) *LoadTester {
	return &LoadTester{
		endpoint: end,
		request:  r,
		client:   &http.Client{},
		conns:    conns,
		dur:      dur,
		rate:     rate,
		stats:    &Stats{Endpoint: end},
	}
}

// Run initializes the LoadTester with its # of conns for a given duration
// passes the results to the statistics channel
func (l *LoadTester) Run(ch chan Stats) {
    var wg sync.WaitGroup

    for i := 0; i < l.conns; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            ticker := time.NewTicker(l.rate)
            defer ticker.Stop()

            stop := time.After(l.dur)
            for {
                select {
                case <-stop:
                    return
                case <-ticker.C:
                    l.test()
                }
            }
        }()
    }
    wg.Wait()

    l.stats.avg()
    ch <- *l.stats
}

func (l *LoadTester) test() {
    var body []byte

    start := time.Now();
    redisDB.IncrementRequestCount();

    resp, err := l.client.Do(l.request);
    rd := time.Since(start);

    if err != nil {
        log.Printf("[ERROR]: Request failed: %v", err)
        l.stats.update(0, 0, err)
        return
    }

    defer resp.Body.Close();

    body, err = io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("[ERROR]: Failed to read response body: %v", err)
        l.stats.update(0, 0, err)
        return
    }

    rs := len(body)
    l.stats.update(rs, rd, nil)
    redisDB.DecrementRequestCount();
}