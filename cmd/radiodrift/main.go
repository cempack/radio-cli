package main

import (
	"fmt"
	"os"

	"github.com/cempack/radio-cli/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := app.New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
