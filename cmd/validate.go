package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Naganathan05/Load-Pulse/utils"
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
		utils.LogError("testConfig configuration is invalid: host must be a non-empty string")
		os.Exit(1)
	}

	if conf.Duration <= 0 {
		utils.LogError("testConfig configuration is invalid: duration must be a positive number")
		os.Exit(1)
	}

	if len(conf.Req) == 0 {
		utils.LogError("testConfig configuration is invalid: requests array must contain at least one request")
		os.Exit(1)
	}

	for i, r := range conf.Req {
		if r.Method == "" {
			utils.LogError(fmt.Sprintf("testConfig configuration is invalid: request[%d].method must be a non-empty string", i))
			os.Exit(1)
		}
		if r.Endpoint == "" {
			utils.LogError(fmt.Sprintf("testConfig configuration is invalid: request[%d].endpoint must be a non-empty string", i))
			os.Exit(1)
		}
		if r.Connections <= 0 {
			utils.LogError(fmt.Sprintf("testConfig configuration is invalid: request[%d].connections must be a positive integer", i))
			os.Exit(1)
		}
		if r.Rate <= 0 {
			utils.LogError(fmt.Sprintf("testConfig configuration is invalid: request[%d].rate must be a positive number", i))
			os.Exit(1)
		}
		if r.ConcurrencyLimit <= 0 {
			utils.LogError(fmt.Sprintf("testConfig configuration is invalid: request[%d].concurrencyLimit must be a positive integer", i))
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

		utils.LogInfo("Validating test configuration file: " + testConfigPath)

		if _, err := os.Stat(testConfigPath); err != nil {
			utils.LogError("testConfig file not found or not accessible: " + err.Error())
			os.Exit(1)
		}

		conf, err := loadTestConfig(testConfigPath)
		if err != nil {
			utils.LogError("Invalid JSON or structure in testConfig file: " + err.Error())
			os.Exit(1)
		}
		validateTestConfig(conf)
		fmt.Printf("testConfig configuration is valid: %s\n", testConfigPath)
	},
}
