package utils

import (
	"os/exec"
)

func IsDockerRunning() (bool, error) {
	cmd := exec.Command("docker", "stats", "--no-stream")
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
