package Tester

import (
	"fmt"
	"time"
	"sync"
)

// Stats holds individual statistics from requests
type Stats struct {
    sync.Mutex
    Endpoint      string
    ResponseSize  float64
    ResponseDur   time.Duration
    TotalRequests int64
    Err           int64
}

// update the object with incoming data
func (s *Stats) update(rs int, rd time.Duration, err error) {
    s.Lock()
    defer s.Unlock()

    s.TotalRequests++
    if err != nil {
        s.Err++
        return
    }
    s.ResponseSize += float64(rs)
    s.ResponseDur += rd
}

// change responseSize, responseDur to be averages based on total requests
// should be called after the tester is done
// does not count errors towards time or size averages
func (s *Stats) avg() {
    s.Lock()
    defer s.Unlock()

    completed := s.TotalRequests - s.Err
    if completed > 0 {
        s.ResponseDur /= time.Duration(completed)
        s.ResponseSize /= float64(completed)
    } else {
        s.ResponseDur = 0
        s.ResponseSize = 0
    }
}

func (s *Stats) print() {
	fmt.Printf("Test completed for endpoint: %s \n", s.Endpoint)
	fmt.Printf("	Total requests completed: %d \n", s.TotalRequests)
	// fmt.Printf("	Total errors: %d \n", s.err)
	fmt.Printf("	Average response size: %f bytes\n", s.ResponseSize)
	fmt.Printf("	Average response time: %s \n", s.ResponseDur.String())
}
