package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	logo = `
  _                      _       ____        _
 | |    ___   __ _  __| |     |  _ \ _   _| |___  ___
 | |   / _ \ / _` + "`" + ` |/ _` + "`" + ` |_____| |_) | | | | / __|/ _ \
 | |__| (_) | (_| | (_| |_____|  __/| |_| | \__ \  __/
 |_____\___/ \__,_|\__,_|     |_|    \__,_|_|___/\___|
`
)

func PrintBanner() {
    // Check if CI environment variable is set
    if os.Getenv("CI") != "" {
        return
    }

    // Check if output is a terminal
    if !term.IsTerminal(int(os.Stdout.Fd())) {
        return
    }

    style := lipgloss.NewStyle().
        Foreground(lipgloss.Color("86")). // Cyan
        Bold(true)

    fmt.Println(style.Render(logo))
}
