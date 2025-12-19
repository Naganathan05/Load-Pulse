package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost             = "http://host.docker.internal:8081/"
	defaultDuration         = 10
	defaultData             = ""
	defaultConnections      = 10
	defaultRate             = 1
	defaultConcurrencyLimit = 1
)

func runTestConfigInitWizard() (*testConfig, error) {
	reader := bufio.NewReader(os.Stdin)

	ColorPrompt("Host [" + defaultHost + "] (press Enter to accept): ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)
	if host == "" {
		host = defaultHost
	}

	ColorPrompt("Duration in seconds [" + strconv.Itoa(defaultDuration) + "]: ")
	durationStr, _ := reader.ReadString('\n')
	durationStr = strings.TrimSpace(durationStr)
	duration := defaultDuration
	if durationStr != "" {
		if v, err := strconv.Atoi(durationStr); err == nil && v > 0 {
			duration = v
		} else {
			return nil, fmt.Errorf("invalid duration value")
		}
	}

	ColorPrompt("Number of request definitions [1]: ")
	reqCountStr, _ := reader.ReadString('\n')
	reqCountStr = strings.TrimSpace(reqCountStr)
	reqCount := 1
	if reqCountStr != "" {
		if v, err := strconv.Atoi(reqCountStr); err == nil && v > 0 {
			reqCount = v
		} else {
			return nil, fmt.Errorf("invalid number of requests")
		}
	}

	var requests []testConfigRequest

	for i := 0; i < reqCount; i++ {
		fmt.Print("\n")
		ColorPrompt(fmt.Sprintf("Configuring request #%d", i+1))
		fmt.Print("\n")

		var method string
		for {
			ColorPrompt("Method (e.g. GET, POST): ")
			method, _ = reader.ReadString('\n')
			method = strings.TrimSpace(method)
			if method != "" {
				break
			}
			ColorHelp("Method is required. Please enter a method.")
		}

		var endpoint string
		for {
			ColorPrompt("Endpoint (e.g. api/admin/getAllDepartments): ")
			endpoint, _ = reader.ReadString('\n')
			endpoint = strings.TrimSpace(endpoint)
			if endpoint != "" {
				break
			}
			ColorHelp("Endpoint is required. Please enter an endpoint.")
		}

		ColorPrompt("Body data (string) [" + defaultData + "]: ")
		data, _ := reader.ReadString('\n')
		data = strings.TrimSpace(data)
		if data == "" {
			data = defaultData
		}

		ColorPrompt("Connections [" + strconv.Itoa(defaultConnections) + "]: ")
		connsStr, _ := reader.ReadString('\n')
		connsStr = strings.TrimSpace(connsStr)
		connections := defaultConnections
		if connsStr != "" {
			if v, err := strconv.Atoi(connsStr); err == nil && v > 0 {
				connections = v
			} else {
				return nil, fmt.Errorf("invalid connections value")
			}
		}

		ColorPrompt("Rate (milliseconds between requests) [" + strconv.Itoa(defaultRate) + "]: ")
		rateStr, _ := reader.ReadString('\n')
		rateStr = strings.TrimSpace(rateStr)
		rate := defaultRate
		if rateStr != "" {
			if v, err := strconv.Atoi(rateStr); err == nil && v > 0 {
				rate = v
			} else {
				return nil, fmt.Errorf("invalid rate value")
			}
		}

		ColorPrompt("Concurrency limit [" + strconv.Itoa(defaultConcurrencyLimit) + "]: ")
		clStr, _ := reader.ReadString('\n')
		clStr = strings.TrimSpace(clStr)
		concurrencyLimit := defaultConcurrencyLimit
		if clStr != "" {
			if v, err := strconv.Atoi(clStr); err == nil && v > 0 {
				concurrencyLimit = v
			} else {
				return nil, fmt.Errorf("invalid concurrency limit value")
			}
		}

		requests = append(requests, testConfigRequest{
			Method:           method,
			Endpoint:         endpoint,
			Data:             data,
			Connections:      connections,
			Rate:             time.Duration(rate),
			ConcurrencyLimit: concurrencyLimit,
		})
	}

	cfg := &testConfig{
		Req:      requests,
		Host:     host,
		Duration: time.Duration(duration),
	}
	validateTestConfig(cfg)
	return cfg, nil
}
