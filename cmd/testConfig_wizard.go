package cmd

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

func runTestConfigInitWizard() (*testConfig, error) {
    reader := bufio.NewReader(os.Stdin)

    // Defaults based on current testConfig.json
    defaultHost := "http://host.docker.internal:8081/"
    defaultDuration := 10
    defaultData := ""
    defaultConnections := 10
    defaultRate := 1
    defaultConcurrencyLimit := 1

    fmt.Printf(
        "%s %s ",
        ColorPrompt("Host"),
        ColorPrompt("["+defaultHost+"] (press Enter to accept):"),
    )
    host, _ := reader.ReadString('\n')
    host = strings.TrimSpace(host)
    if host == "" {
        host = defaultHost
    }

    fmt.Printf(
        "%s %s ",
        ColorPrompt("Duration in seconds"),
        ColorPrompt("["+strconv.Itoa(defaultDuration)+"]:"+""), // e.g. [10]:
    )
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

    fmt.Printf(
        "%s %s ",
        ColorPrompt("Number of request definitions"),
        ColorPrompt("[1]:"+""), // default 1
    )
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
        fmt.Printf(
            "\n%s\n",
            ColorPrompt(fmt.Sprintf("Configuring request #%d", i+1)),
        )

        var method string
        for {
            fmt.Printf("%s ", ColorPrompt("Method (e.g. GET, POST):"))
            method, _ = reader.ReadString('\n')
            method = strings.TrimSpace(method)
            if method != "" {
                break
            }
            fmt.Println(ColorHelp("Method is required. Please enter a method."))
        }

        var endpoint string
        for {
            fmt.Printf("%s ", ColorPrompt("Endpoint (e.g. api/admin/getAllDepartments):"))
            endpoint, _ = reader.ReadString('\n')
            endpoint = strings.TrimSpace(endpoint)
            if endpoint != "" {
                break
            }
            fmt.Println(ColorHelp("Endpoint is required. Please enter an endpoint."))
        }

        fmt.Printf(
            "%s %s ",
            ColorPrompt("Body data (string)"),
            ColorPrompt("["+defaultData+"]:"+""), // default empty
        )
        data, _ := reader.ReadString('\n')
        data = strings.TrimSpace(data)
        if data == "" {
            data = defaultData
        }

        fmt.Printf(
            "%s %s ",
            ColorPrompt("Connections"),
            ColorPrompt("["+strconv.Itoa(defaultConnections)+"]:"+""), // e.g. [10]:
        )
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

        fmt.Printf(
            "%s %s ",
            ColorPrompt("Rate (milliseconds between requests)"),
            ColorPrompt("["+strconv.Itoa(defaultRate)+"]:"+""), // e.g. [1]:
        )
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

        fmt.Printf(
            "%s %s ",
            ColorPrompt("Concurrency limit"),
            ColorPrompt("["+strconv.Itoa(defaultConcurrencyLimit)+"]:"+""), // e.g. [1]:
        )
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

    if err := validateTestConfig(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}