package utils

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func LogInfo(msg string) {
	fmt.Println(infoStyle.Render("[LOG]: " + msg));
}

func LogError(msg string) {
	fmt.Println(errorStyle.Render("[ERROR]: " + msg));
}