package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "loadpulse",
    Short: "LoadPulse - Load Testing Tool",
	Long: `A CLI tool for Load Testing Web Applications.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Print(configs.LogoStyle.Render(logo.LOGO_ASCII));
	// },
}

func Execute() {
    PrintBanner();
    cobra.CheckErr(rootCmd.Execute());
}