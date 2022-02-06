package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 3)

var WindowDimensions = [2]int{0, 0}
var State string = "main"

type model struct {
	menus map[string]tea.Model
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
		menus: map[string]tea.Model{
			"main":  initMainScreen(),
			"shell": initialShellScreen(),
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()

		WindowDimensions[0] = msg.Width - left - right
		WindowDimensions[1] = msg.Height - top - bottom - 2

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "backspace", "left", "h":
			State = "main"
			return m, nil
		}
	}

	_, cmd := m.menus[State].Update(msg)

	return m, cmd
}

func (m model) View() string {
	s := ""

	// If there's a menu to render, render it
	if State != "" {
		s += m.menus[State].View()
	}

	// The footer
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7d8899"))
	s += style.Render("\nq / Ctrl + C to quit\n")

	return docStyle.Render(s)
}
