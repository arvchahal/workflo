package main

import (
	"fmt"
	"os"

	"workflo/cli"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(cli.NewModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
