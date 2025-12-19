package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // blue
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // magenta
	infoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
    errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red
)

func LogPrompt(text string) {
	fmt.Print(promptStyle.Render(text))
}

func LogHelp(text string) {
	fmt.Println(helpStyle.Render(text))
}

func LogInfo(msg string) {
    fmt.Println(infoStyle.Render("[LOG]: " + msg))
}

func LogError(msg string) {
    fmt.Println(errorStyle.Render("[ERROR]: " + msg))
}

func LogPlain(text string) {
	fmt.Print(text)
}

func LogNewLine() {
	fmt.Println()
}