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

func validateTestConfig(conf *testConfig) error {
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

        if err := validateTestConfig(conf); err != nil {
            utils.LogError("testConfig configuration is invalid: " + err.Error())
            os.Exit(1)
        }

        fmt.Printf("testConfig configuration is valid: %s\n", testConfigPath)
    },
}

func init() {
    rootCmd.AddCommand(validateCmd)
}