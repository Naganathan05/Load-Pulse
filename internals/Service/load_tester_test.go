package Service

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSinglePostRequest(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if r.Method == "POST" && len(body) == 0 {
			t.Logf("Received empty POST body")
			w.WriteHeader(http.StatusBadRequest) // Fail if empty
			return
		}
		t.Logf("Received POST body: %s", string(body))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a request with body
	data := "test payload"
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBufferString(data))

	// Create LoadTester instance
	lt := NewTester(req, 1, time.Second, time.Millisecond, server.URL, 1, []byte(data))

	// First request
	t.Log("Sending First Request")
	resp1, err := lt.DoRequest()
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}
	resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("First request status: %d", resp1.StatusCode)
	}

	// Second request (simulate reuse)
	t.Log("Sending Second Request")
	resp2, err := lt.DoRequest()
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Second request status: %d. Expected 200, got %d. Body was likely empty.", resp2.StatusCode, resp2.StatusCode)
	}
}

func TestConcurrentPostRequests(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if r.Method == "POST" && len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	data := "concurrent payload"
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBufferString(data))
	lt := NewTester(req, 10, time.Second, time.Millisecond, server.URL, 1, []byte(data))

	// Run concurrently
	concurrency := 10
	errCh := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < 5; j++ {
				resp, err := lt.DoRequest()
				if err != nil {
					errCh <- err
					return
				}
				resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					errCh <- io.EOF // represent empty body error
					return
				}
			}
			errCh <- nil
		}()
	}

	for i := 0; i < concurrency; i++ {
		err := <-errCh
		if err != nil {
			t.Fatalf("Concurrent request failed: %v", err)
		}
	}
}
