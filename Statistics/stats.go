package Statistics

import (
	"fmt"
	"time"
	"sync"
)

type Stats struct {
    sync.Mutex
    Endpoint      string
    ResponseSize  float64
    ResponseDur   time.Duration
    FailedRequests int64
    TotalRequests int64
}

func (s *Stats) Update(rs int, rd time.Duration, err error) {
    s.Lock()
    defer s.Unlock()

    s.TotalRequests++
    if err != nil {
        s.FailedRequests++
        return
    }
    s.ResponseSize += float64(rs)
    s.ResponseDur += rd
}

func (s *Stats) Avg() {
    s.Lock()
    defer s.Unlock()

    completed := s.TotalRequests - s.FailedRequests
    if completed > 0 {
        s.ResponseDur /= time.Duration(completed)
        s.ResponseSize /= float64(completed)
    } else {
        s.ResponseDur = 0
        s.ResponseSize = 0
    }
}

func (s *Stats) Print() {
	fmt.Printf("\n[LOG]: Test completed for endpoint: %s \n", s.Endpoint)
    fmt.Println("----------------------- STATS -------------------------------");
	fmt.Printf("	Total requests completed: %d \n", s.TotalRequests)
	fmt.Printf("	Average response size: %f bytes\n", s.ResponseSize)
	fmt.Printf("	Average response time: %s \n", s.ResponseDur.String())
    fmt.Printf("--------------------------------------------------------------\n\n");
}