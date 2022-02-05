package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	menus map[string]Screen
	state string
}

func Init() {
	program := tea.NewProgram(initialModel())
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "backspace":
			m.state = "main"
			return m, nil
		}
	}

	cmd := m.menus[m.state].Update(msg, &m)

	return m, cmd
}

func (m model) View() string {
	// The header
	s := "What do you want to do?\n\n"

	// If there's a menu to render, render it
	if m.state != "" {
		s += m.menus[m.state].Render()
	}

	// The footer
	s += "\nPress q to quit.\n"

	return s
}
