package utils

import (
    "fmt"

    "github.com/charmbracelet/lipgloss"
)

var (
    infoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
    errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red
    promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // blue
    helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // magenta
)

func LogInfo(msg string) {
    fmt.Println(infoStyle.Render("[LOG]: " + msg))
}

func LogError(msg string) {
    fmt.Println(errorStyle.Render("[ERROR]: " + msg))
}

func ColorPrompt(text string) string {
    return promptStyle.Render(text)
}

func ColorHelp(text string) string {
    return helpStyle.Render(text)
}