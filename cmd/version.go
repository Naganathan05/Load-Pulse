package cmd

import (
	"fmt"
	"Load-Pulse/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of LoadPulse",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("LoadPulse version:", utils.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}