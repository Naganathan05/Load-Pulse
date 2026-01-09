package Service

import (
	"io"
	"time"
	"bytes"
	"net/http"

	"Load-Pulse/Statistics"
)

type Bench struct {
	Testers []*LoadTester
	Ch      chan *Statistics.Stats
}

type LoadTester struct {
	Endpoint         string
	Conns            int
	Request          *http.Request
	Client           *http.Client
	Stats            *Statistics.Stats
	Dur              time.Duration
	Rate             time.Duration
	ConcurrencyLimit int
	RequestBody      []byte
}

func Min(a int, b int) int {
	if a < b {
		return a;
	}
	return b;
}

func Max(a int, b int) int {
	if a > b {
		return a;
	}
	return b;
}

func NewTester(r *http.Request, conns int, dur, rate time.Duration, end string, concurrencyLimit int, reqBody []byte) *LoadTester {
	return &LoadTester{
		Endpoint:         end,
		Request:          r,
		Client:           &http.Client{},
		Conns:            conns,
		Dur:              dur,
		Rate:             rate,
		Stats:            &Statistics.Stats{Endpoint: end},
		ConcurrencyLimit: concurrencyLimit,
		RequestBody:      reqBody,
	}
}

func NewLoadTester(path string) (*Bench, error) {
	var testers []*LoadTester;

	conf, err := FromJSON(path);
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

		lt := NewTester(r, req.Connections, conf.Duration*time.Second, req.Rate*time.Millisecond, req.Endpoint, req.ConcurrencyLimit, []byte(req.Data))
		testers = append(testers, lt);
	}

	b := &Bench{
		Testers: testers,
		Ch:      make(chan *Statistics.Stats, len(testers)),
	}

	return b, nil;
}

func (l *LoadTester) DoRequest() (*http.Response, error) {
	req := l.Request.Clone(l.Request.Context())
	if len(l.RequestBody) > 0 {
		req.Body = io.NopCloser(bytes.NewReader(l.RequestBody))
		req.ContentLength = int64(len(l.RequestBody))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(l.RequestBody)), nil
		}
	}
	return l.Client.Do(req)
}