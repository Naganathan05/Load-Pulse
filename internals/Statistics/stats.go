package Statistics

import (
	"fmt"
	"time"
	"sync"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	White  = "\033[37m"
	Reset  = "\033[0m"
	Violet    = "\033[35m"
)

func LogServer(message string) {
	fmt.Printf("%s%s%s", White, message, Reset);
}

func LogBlue(message string) {
	fmt.Printf("%s%s%s", Blue, message, Reset);
}

func LogGreen(message string) {
	fmt.Printf("%s%s%s", Green, message, Reset);
}

type Stats struct {
    sync.Mutex
    Endpoint      string
    ResponseSize  float64
    ResponseDur   time.Duration
    FailedRequests int64
    TotalRequests int64
    MaxResponseTime time.Duration
    MinResponseTime time.Duration
}

func (s *Stats) Update(rs int, rd time.Duration, err error) {
    s.Lock();
    defer s.Unlock();

    s.TotalRequests += 1;
    if err != nil {
        s.FailedRequests += 1;
        return;
    }
    s.ResponseSize += float64(rs);
    s.ResponseDur += rd;
    s.MaxResponseTime = max(s.MaxResponseTime, rd);
    s.MinResponseTime = min(s.MinResponseTime, rd);
}

func (s *Stats) Avg() {
    s.Lock();
    defer s.Unlock();

    completed := s.TotalRequests - s.FailedRequests;
    if completed > 0 {
        s.ResponseDur /= time.Duration(completed);
        s.ResponseSize /= float64(completed);
    } else {
        s.ResponseDur = 0;
        s.ResponseSize = 0;
    }
}

func (s *Stats) Print() {
	logMsg := fmt.Sprintf("\n[LOG]: Test completed for endpoint: %s \n", s.Endpoint); LogServer(logMsg);
    LogBlue("------------------------------ STATS --------------------------------\n");
	logMsg = fmt.Sprintf("[STATS]: Total requests completed: %d\n", s.TotalRequests); LogGreen(logMsg);
    logMsg = fmt.Sprintf("[STATS]: Total Number of Error Requests: %d\n", s.FailedRequests); LogGreen(logMsg);
	logMsg = fmt.Sprintf("[STATS]: Average response size: %f bytes\n", s.ResponseSize); LogGreen(logMsg);
    logMsg = fmt.Sprintf("[STATS]: Min response time: %s \n", s.MinResponseTime.String()); LogGreen(logMsg);
    logMsg = fmt.Sprintf("[STATS]: Max response time: %s \n", s.MaxResponseTime.String()); LogGreen(logMsg);
	logMsg = fmt.Sprintf("[STATS]: Average response time: %s \n", s.ResponseDur.String()); LogGreen(logMsg);
    LogBlue("---------------------------------------------------------------------\n\n");
}