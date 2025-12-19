package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Naganathan05/Load-Pulse/utils"
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

		ColorPrompt("Enter output JSON file path [" + defaultPath + "]: ")
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

		ColorPrompt("This wizard will create a testConfig configuration file at: ")
		ColorPrompt(testConfigInitOutputPath)
		fmt.Println()

		cfg, err := runTestConfigInitWizard()
		if err != nil {
			utils.LogError("Failed to create testConfig configuration: " + err.Error())
			os.Exit(1)
		}

		data, err := json.MarshalIndent(cfg, "", "    ")
		if err != nil {
			utils.LogError("Failed to encode testConfig configuration as JSON: " + err.Error())
			os.Exit(1)
		}

		if err := os.WriteFile(testConfigInitOutputPath, data, 0644); err != nil {
			utils.LogError("Failed to write testConfig configuration file: " + err.Error())
			os.Exit(1)
		}

		ColorPrompt("testConfig configuration written to: ")
		ColorPrompt(testConfigInitOutputPath)
		fmt.Println()
		ColorHelp("You can validate it with:")
		fmt.Print("  ")
		ColorHelp("go run .\\main.go validate " + testConfigInitOutputPath)
	},
}
