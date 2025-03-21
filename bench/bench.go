package bench

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Bench is a struct that controls the testers, channel communication
// and stat aggregation
type Bench struct {
	testers []*LoadTester
	ch      chan Stats
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
		// init new Tester with given request
		lt := NewTester(r, req.Connections, conf.Duration*time.Second, req.Rate*time.Millisecond, req.Endpoint)
		testers = append(testers, lt)
	}

	b := &Bench{
		testers: testers,
		ch:      make(chan Stats, len(testers)),
	}

	return b, nil
}

// Run a benchmark test with given config
// run each tester concurrently and wait for them to finish
func (b *Bench) Run() {
	var wg sync.WaitGroup

	for _, tester := range b.testers {
		wg.Add(1)

		go func(t *LoadTester) {
			defer wg.Done()
			// run loadtester with the specific channel
			fmt.Printf("Running test on %s with %d connections making a request every %s \n", t.endpoint, t.conns, t.rate.String())
			t.Run(b.ch)
		}(tester)

	}
	wg.Wait()
	close(b.ch)

	for stat := range b.ch {
		stat.print()
	}
}
