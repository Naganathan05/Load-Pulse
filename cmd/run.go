package cmd

import (
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/Naganathan05/Load-Pulse/utils"
	
)
var verbose bool

var testConfigPath string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the load testing tool",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
	        LogInfo("Initializing Load Pulse")
	        LogInfo("Docker compose command started")
}

		LogInfo("Using testConfig configuration file: " + testConfigPath)
		ok, _ := utils.IsDockerRunning()
		if !ok {
			LogError("Docker Engine Not Running. Please Start Docker Daemon and try again.\n")
			os.Exit(1)
		}

		LogInfo("Spinning up Docker Containers...")
		startCmd := exec.Command("docker", "compose", "up", "-d", "--build")

		env := os.Environ()
		env = append(env, "TESTCONFIG_FILE_PATH="+testConfigPath)
		startCmd.Env = env

		/*-------------------------- Debugging --------------------------
		startCmd.Stdout = os.Stdout;
		startCmd.Stderr = os.Stderr;
		---------------------------------------------------------------*/
		startCmd.Stdout = nil
		startCmd.Stderr = nil

		if err := startCmd.Run(); err != nil {
			LogError("Failed to start containers with Docker Compose: " + err.Error())
			os.Exit(1)
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = "Load Testing in Progress..."
		s.Start()
		for {
			out, _ := exec.Command("docker", "inspect", "--format", "{{.State.Running}}", "aggregator").Output()
			if string(out) == "false\n" {
				s.Stop()
				break
			}
			time.Sleep(2 * time.Second)
		}

		LogInfo("Load Test Completed. Logging the Aggregator Container Logs: ")
		logsCmd := exec.Command("docker", "logs", "aggregator")
		logsCmd.Stdout = os.Stdout
		logsCmd.Stderr = os.Stderr
		logsCmd.Run()

		cleanCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	runCmd.Flags().StringVarP(
		&testConfigPath,
		"config",
		"c",
		"testConfig.json",
		"Path to testConfig configuration file",
	)
     runCmd.Flags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"Enable verbose logging",
	)
}