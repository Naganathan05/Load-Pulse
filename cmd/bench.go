package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Naganathan05/Load-Pulse/utils"
	"github.com/spf13/cobra"
)

type benchConfig struct {
	Req      []benchRequest `json:"requests"`
	Host     string         `json:"host"`
	Duration time.Duration  `json:"duration"`
}

type benchRequest struct {
	Method           string        `json:"method"`
	Endpoint         string        `json:"endpoint"`
	Data             string        `json:"data"`
	Connections      int           `json:"connections"`
	Rate             time.Duration `json:"rate"`
	ConcurrencyLimit int           `json:"concurrencyLimit"`
}

func loadBenchConfig(path string) (*benchConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var conf benchConfig
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Bench configuration related commands",
}

var benchValidateCmd = &cobra.Command{
	Use:   "validate [bench-file]",
	Short: "Validate a bench configuration JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		benchPath := "bench.json"
		if len(args) == 1 {
			benchPath = args[0]
		}

		utils.LogInfo("Validating bench configuration file: " + benchPath)

		if _, err := os.Stat(benchPath); err != nil {
			utils.LogError("Bench file not found or not accessible: " + err.Error())
			os.Exit(1)
		}

		conf, err := loadBenchConfig(benchPath)
		if err != nil {
			utils.LogError("Invalid JSON or structure in bench file: " + err.Error())
			os.Exit(1)
		}

		if err := validateBenchConfig(conf); err != nil {
			utils.LogError("Bench configuration is invalid: " + err.Error())
			os.Exit(1)
		}

		fmt.Printf("Bench configuration is valid: %s\n", benchPath)
	},
}

var benchInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively create a bench configuration JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Default output path: bench.json in current directory
		defaultPath := "bench.json"

		fmt.Printf(
			"%s %s ",
			utils.ColorPrompt("Enter output JSON file path"),
			utils.ColorPrompt("["+defaultPath+"]:"),
		)
		pathInput, _ := reader.ReadString('\n')
		pathInput = strings.TrimSpace(pathInput)
		if pathInput == "" {
			pathInput = defaultPath
		} else {
			// Ensure the file has a .json extension
			if !strings.HasSuffix(strings.ToLower(pathInput), ".json") {
				pathInput = pathInput + ".json"
			}
		}
		benchInitOutputPath := pathInput

		fmt.Printf(
			"%s %s\n",
			utils.ColorPrompt("This wizard will create a bench configuration file at:"),
			utils.ColorPrompt(benchInitOutputPath),
		)

		cfg, err := runBenchInitWizard()
		if err != nil {
			utils.LogError("Failed to create bench configuration: " + err.Error())
			os.Exit(1)
		}

		data, err := json.MarshalIndent(cfg, "", "    ")
		if err != nil {
			utils.LogError("Failed to encode bench configuration as JSON: " + err.Error())
			os.Exit(1)
		}

		if err := os.WriteFile(benchInitOutputPath, data, 0644); err != nil {
			utils.LogError("Failed to write bench configuration file: " + err.Error())
			os.Exit(1)
		}

		fmt.Printf(
			"%s %s\n",
			utils.ColorPrompt("Bench configuration written to:"),
			utils.ColorPrompt(benchInitOutputPath),
		)
		fmt.Println(utils.ColorHelp("You can validate it with:"))
		fmt.Printf(
			"  %s\n",
			utils.ColorHelp("go run .\\main.go bench validate "+benchInitOutputPath),
		)
	},
}

func init() {
	rootCmd.AddCommand(benchCmd)
	benchCmd.AddCommand(benchValidateCmd)
	benchCmd.AddCommand(benchInitCmd)
}

func validateBenchConfig(conf *benchConfig) error {
	if conf.Host == "" {
		return fmt.Errorf("host must be a non-empty string")
	}

	if conf.Duration <= 0 {
		return fmt.Errorf("duration must be a positive number")
	}

	if len(conf.Req) == 0 {
		return fmt.Errorf("requests array must contain at least one request")
	}

	for i, r := range conf.Req {
		if r.Method == "" {
			return fmt.Errorf("request[%d].method must be a non-empty string", i)
		}
		if r.Endpoint == "" {
			return fmt.Errorf("request[%d].endpoint must be a non-empty string", i)
		}
		if r.Connections <= 0 {
			return fmt.Errorf("request[%d].connections must be a positive integer", i)
		}
		if r.Rate <= 0 {
			return fmt.Errorf("request[%d].rate must be a positive number", i)
		}
		if r.ConcurrencyLimit <= 0 {
			return fmt.Errorf("request[%d].concurrencyLimit must be a positive integer", i)
		}
		// r.Data can be empty.
	}

	return nil
}

func runBenchInitWizard() (*benchConfig, error) {
	reader := bufio.NewReader(os.Stdin)

	// Defaults based on your current bench.json
	defaultHost := "http://host.docker.internal:8081/"
	defaultDuration := 10
	defaultData := ""
	defaultConnections := 10
	defaultRate := 1
	defaultConcurrencyLimit := 1

	fmt.Printf(
		"%s %s ",
		utils.ColorPrompt("Host"),
		utils.ColorPrompt("["+defaultHost+"] (press Enter to accept):"),
	)
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)
	if host == "" {
		host = defaultHost
	}

	fmt.Printf(
		"%s %s ",
		utils.ColorPrompt("Duration in seconds"),
		utils.ColorPrompt("["+strconv.Itoa(defaultDuration)+"]:"+""), // e.g. [10]:
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
		utils.ColorPrompt("Number of request definitions"),
		utils.ColorPrompt("[1]:"+""), // default 1
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

	var requests []benchRequest

	for i := 0; i < reqCount; i++ {
		fmt.Printf(
			"\n%s\n",
			utils.ColorPrompt(fmt.Sprintf("Configuring request #%d", i+1)),
		)

		var method string
		for {
			fmt.Printf("%s ", utils.ColorPrompt("Method (e.g. GET, POST):"))
			method, _ = reader.ReadString('\n')
			method = strings.TrimSpace(method)
			if method != "" {
				break
			}
			fmt.Println(utils.ColorHelp("Method is required. Please enter a method."))
		}

		var endpoint string
		for {
			fmt.Printf("%s ", utils.ColorPrompt("Endpoint (e.g. api/admin/getAllDepartments):"))
			endpoint, _ = reader.ReadString('\n')
			endpoint = strings.TrimSpace(endpoint)
			if endpoint != "" {
				break
			}
			fmt.Println(utils.ColorHelp("Endpoint is required. Please enter an endpoint."))
		}

		fmt.Printf(
			"%s %s ",
			utils.ColorPrompt("Body data (string)"),
			utils.ColorPrompt("["+defaultData+"]:"+""), // default empty
		)
		data, _ := reader.ReadString('\n')
		data = strings.TrimSpace(data)
		if data == "" {
			data = defaultData
		}

		fmt.Printf(
			"%s %s ",
			utils.ColorPrompt("Connections"),
			utils.ColorPrompt("["+strconv.Itoa(defaultConnections)+"]:"+""), // e.g. [10]:
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
			utils.ColorPrompt("Rate (milliseconds between requests)"),
			utils.ColorPrompt("["+strconv.Itoa(defaultRate)+"]:"+""), // e.g. [1]:
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
			utils.ColorPrompt("Concurrency limit"),
			utils.ColorPrompt("["+strconv.Itoa(defaultConcurrencyLimit)+"]:"+""), // e.g. [1]:
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

		requests = append(requests, benchRequest{
			Method:           method,
			Endpoint:         endpoint,
			Data:             data,
			Connections:      connections,
			Rate:             time.Duration(rate),
			ConcurrencyLimit: concurrencyLimit,
		})
	}

	cfg := &benchConfig{
		Req:      requests,
		Host:     host,
		Duration: time.Duration(duration),
	}

	if err := validateBenchConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
