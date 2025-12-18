package cmd

import (
    "github.com/charmbracelet/lipgloss"
)

var (
    promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // blue
    helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // magenta
)

func ColorPrompt(text string) string {
    return promptStyle.Render(text)
}

func ColorHelp(text string) string {
    return helpStyle.Render(text)
}