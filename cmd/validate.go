package cmd

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"github.com/spf13/cobra"
)

func loadTestConfig(path string) (*testConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var conf testConfig
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func validateTestConfig(conf *testConfig) {
	if conf.Host == "" {
		LogError("testConfig configuration is invalid: host must be a non-empty string")
		os.Exit(1)
	}

	if conf.Duration <= 0 {
		LogError("testConfig configuration is invalid: duration must be a positive number")
		os.Exit(1)
	}

	if len(conf.Req) == 0 {
		LogError("testConfig configuration is invalid: requests array must contain at least one request")
		os.Exit(1)
	}

	for i, r := range conf.Req {
		if r.Method == "" {
			LogError("testConfig configuration is invalid: request[" + strconv.Itoa(i) + "].method must be a non-empty string")
			os.Exit(1)
		}
		if r.Endpoint == "" {
			LogError("testConfig configuration is invalid: request[" + strconv.Itoa(i) + "].endpoint must be a non-empty string")
			os.Exit(1)
		}
		if r.Connections <= 0 {
			LogError("testConfig configuration is invalid: request[" + strconv.Itoa(i) + "].connections must be a positive integer")
			os.Exit(1)
		}
		if r.Rate <= 0 {
			LogError("testConfig configuration is invalid: request[" + strconv.Itoa(i) + "].rate must be a positive number")
			os.Exit(1)
		}
		if r.ConcurrencyLimit <= 0 {
			LogError("testConfig configuration is invalid: request[" + strconv.Itoa(i) + "].concurrencyLimit must be a positive integer")
			os.Exit(1)
		}
		// r.Data can be empty.
	}
}

var validateCmd = &cobra.Command{
	Use:   "validate [testConfig-file]",
	Short: "Validate a test configuration JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		testConfigPath := "testConfig.json"
		if len(args) == 1 {
			testConfigPath = args[0]
		}

		LogInfo("Validating test configuration file: " + testConfigPath)

		if _, err := os.Stat(testConfigPath); err != nil {
			LogError("testConfig file not found or not accessible: " + err.Error())
			os.Exit(1)
		}

		conf, err := loadTestConfig(testConfigPath)
		if err != nil {
			LogError("Invalid JSON or structure in testConfig file: " + err.Error())
			os.Exit(1)
		}
		validateTestConfig(conf)
		LogInfo("testConfig configuration is valid: " + testConfigPath)
	},
}
