package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	menus  map[string]Screen
	state  string
	width  int
	height int
}

func Init() {
	program := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := program.Start(); err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		menus: map[string]Screen{
			"main":  initMainScreen(),
			"build": &BuildScreen{},
			"run":   &RunScreen{},
			"shell": &ShellScreen{},
		},
		state: "main",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "backspace", "left", "h":
			m.state = "main"
			return m, nil
		}
	}

	cmd := m.menus[m.state].Update(msg, &m)

	return m, cmd
}

func (m model) View() string {
	s := ""

	// If there's a menu to render, render it
	if m.state != "" {
		s += m.menus[m.state].Render(&m)
	}

	// The footer
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7d8899"))
	s += style.Render("\nq / Ctrl + C to quit\n")

	return s
}
