package cmd

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"
	"github.com/spf13/cobra"
)

type testConfig struct {
	Req      []testConfigRequest `json:"requests"`
	Host     string              `json:"host"`
	Duration time.Duration       `json:"duration"`
}

type testConfigRequest struct {
	Method           string        `json:"method"`
	Endpoint         string        `json:"endpoint"`
	Data             string        `json:"data"`
	Connections      int           `json:"connections"`
	Rate             time.Duration `json:"rate"`
	ConcurrencyLimit int           `json:"concurrencyLimit"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively create a testConfig configuration JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Default output path: testConfig.json in current directory
		defaultPath := "testConfig.json"

		LogPrompt("Enter output JSON file path [" + defaultPath + "]: ")
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
		testConfigInitOutputPath := pathInput

		LogPrompt("This wizard will create a testConfig configuration file at: ")
		LogPrompt(testConfigInitOutputPath)
		LogNewLine()

		cfg, err := runTestConfigInitWizard()
		if err != nil {
			LogError("Failed to create testConfig configuration: " + err.Error())
			os.Exit(1)
		}

		data, err := json.MarshalIndent(cfg, "", "    ")
		if err != nil {
			LogError("Failed to encode testConfig configuration as JSON: " + err.Error())
			os.Exit(1)
		}

		if err := os.WriteFile(testConfigInitOutputPath, data, 0644); err != nil {
			LogError("Failed to write testConfig configuration file: " + err.Error())
			os.Exit(1)
		}

		LogPrompt("testConfig configuration written to: ")
		LogPrompt(testConfigInitOutputPath)
		LogNewLine()
		LogHelp("You can validate it with:")
		LogPlain("  ")
		LogHelp("go run .\\main.go validate " + testConfigInitOutputPath)
	},
}
