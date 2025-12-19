package cmd

import (
	"os"
	"os/exec"
	"github.com/spf13/cobra"
);

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up containers",
	Long:  "Clean up containers and other associated resources",
	Run: func(cmd *cobra.Command, args []string) {
		LogInfo("Starting cleanup process ...");

		out, err := exec.Command("docker", "compose", "ps", "-q").Output();
		if err != nil || len(out) == 0 {
			LogInfo("No containers found to clean up");
			return;
		};

		stopCmd := exec.Command("docker", "compose", "down", "-v");

		if err := stopCmd.Run(); err != nil {
			LogError("Failed to stop containers with Docker Compose: " + err.Error());
			os.Exit(1);
		};

		LogInfo("Container Cleanup Successfully Completed");
	},
};

func init() {
	rootCmd.AddCommand(cleanCmd);
};