package cmd

import (
	"os"
	"fmt"
	"os/exec"
	"Load-Pulse/utils"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"time"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the load testing tool",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Initializing Load Pulse");
		ok, err := utils.IsDockerRunning()
		if !ok {
			fmt.Printf("Docker engine not running: %v", err)
			os.Exit(1)
		}		
		
		utils.LogInfo("Spinning up Docker Containers...");
		startCmd := exec.Command("docker", "compose", "up", "-d", "--build")
		startCmd.Stdout = nil;
		startCmd.Stderr = nil;
		
		if err := startCmd.Run(); err != nil {
			utils.LogError("Failed to start containers with Docker Compose: " + err.Error())
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

		utils.LogInfo("Load Test Completed. Logging the Aggregator Container Logs: ")
		logsCmd := exec.Command("docker", "logs", "aggregator")
		logsCmd.Stdout = os.Stdout
		logsCmd.Stderr = os.Stderr
		logsCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}