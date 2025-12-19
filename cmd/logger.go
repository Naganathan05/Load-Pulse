package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // blue
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // magenta
)

func ColorPrompt(text string) {
	fmt.Print(promptStyle.Render(text))
}

func ColorHelp(text string) {
	fmt.Println(helpStyle.Render(text))
}
