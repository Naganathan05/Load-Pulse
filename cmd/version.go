package cmd

import (
	"fmt"
	"github.com/Naganathan05/Load-Pulse/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of LoadPulse",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Load Pulse version:", utils.Version);
	},
}

func init() {
	rootCmd.AddCommand(versionCmd);
}