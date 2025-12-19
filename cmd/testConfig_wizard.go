package cmd

import (
	"bufio"
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

	LogPrompt("Host [" + defaultHost + "] (press Enter to accept): ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)
	if host == "" {
		host = defaultHost
	}

	LogPrompt("Duration in seconds [" + strconv.Itoa(defaultDuration) + "]: ")
	durationStr, _ := reader.ReadString('\n')
	durationStr = strings.TrimSpace(durationStr)
	duration := defaultDuration
	if durationStr != "" {
		if v, err := strconv.Atoi(durationStr); err == nil && v > 0 {
			duration = v
		} else {
			LogError("invalid duration value")
			os.Exit(1)
		}
	}

	LogPrompt("Number of request definitions [1]: ")
	reqCountStr, _ := reader.ReadString('\n')
	reqCountStr = strings.TrimSpace(reqCountStr)
	reqCount := 1
	if reqCountStr != "" {
		if v, err := strconv.Atoi(reqCountStr); err == nil && v > 0 {
			reqCount = v
		} else {
			LogError("invalid number of requests")
			os.Exit(1)
		}
	}

	var requests []testConfigRequest

	for i := 0; i < reqCount; i++ {
		LogNewLine()
		LogPrompt("Configuring request #" + strconv.Itoa(i+1))
		LogNewLine()

		var method string
		for {
			LogPrompt("Method (e.g. GET, POST): ")
			method, _ = reader.ReadString('\n')
			method = strings.TrimSpace(method)
			if method != "" {
				break
			}
			LogHelp("Method is required. Please enter a method.")
		}

		var endpoint string
		for {
			LogPrompt("Endpoint (e.g. api/admin/getAllDepartments): ")
			endpoint, _ = reader.ReadString('\n')
			endpoint = strings.TrimSpace(endpoint)
			if endpoint != "" {
				break
			}
			LogHelp("Endpoint is required. Please enter an endpoint.")
		}

		LogPrompt("Body data (string) [" + defaultData + "]: ")
		data, _ := reader.ReadString('\n')
		data = strings.TrimSpace(data)
		if data == "" {
			data = defaultData
		}

		LogPrompt("Connections [" + strconv.Itoa(defaultConnections) + "]: ")
		connsStr, _ := reader.ReadString('\n')
		connsStr = strings.TrimSpace(connsStr)
		connections := defaultConnections
		if connsStr != "" {
			if v, err := strconv.Atoi(connsStr); err == nil && v > 0 {
				connections = v
			} else {
				LogError("invalid connections value")
				os.Exit(1)
			}
		}

		LogPrompt("Rate (milliseconds between requests) [" + strconv.Itoa(defaultRate) + "]: ")
		rateStr, _ := reader.ReadString('\n')
		rateStr = strings.TrimSpace(rateStr)
		rate := defaultRate
		if rateStr != "" {
			if v, err := strconv.Atoi(rateStr); err == nil && v > 0 {
				rate = v
			} else {
				LogError("invalid rate value")
				os.Exit(1)
			}
		}

		LogPrompt("Concurrency limit [" + strconv.Itoa(defaultConcurrencyLimit) + "]: ")
		clStr, _ := reader.ReadString('\n')
		clStr = strings.TrimSpace(clStr)
		concurrencyLimit := defaultConcurrencyLimit
		if clStr != "" {
			if v, err := strconv.Atoi(clStr); err == nil && v > 0 {
				concurrencyLimit = v
			} else {
				LogError("invalid concurrency limit value")
				os.Exit(1)
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
