package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	logo = `
██╗      ██████╗  █████╗ ██████╗       ██████╗ ██╗   ██╗██╗     ███████╗███████╗
██║     ██╔═══██╗██╔══██╗██╔══██╗      ██╔══██╗██║   ██║██║     ██╔════╝██╔════╝
██║     ██║   ██║███████║██║  ██║█████╗██████╔╝██║   ██║██║     ███████╗█████╗  
██║     ██║   ██║██╔══██║██║  ██║╚════╝██╔═══╝ ██║   ██║██║     ╚════██║██╔══╝  
███████╗╚██████╔╝██║  ██║██████╔╝      ██║     ╚██████╔╝███████╗███████║███████╗
╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚═════╝       ╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝
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
        Foreground(lipgloss.Color("208")). // Orange
        Bold(true)

    fmt.Println(style.Render(logo))
}
