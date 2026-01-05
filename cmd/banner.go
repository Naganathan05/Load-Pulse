package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	logo = `
    __                    __      ____        __         
   / /   ____  ____ _____/ /_____/ __ \__  __/ /_________
  / /   / __ \/ __ \/ __  /_____/ /_/ / / / / / ___/ _  /
 / /___/ /_/ / /_/ / /_/ /     / ____/ /_/ / (__  )  __/ 
/_____/\____/\__,_/\__,_/     /_/    \__,_/_/____/\___/  
                                                         
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
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4")).
        PaddingTop(1).
        PaddingBottom(1).
        PaddingLeft(4).
        PaddingRight(4).
        MarginBottom(1)

    fmt.Println(style.Render(logo))
}
