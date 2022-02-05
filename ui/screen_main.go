package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainScreen struct {
	choices []string
	cursor  int
}

func initMainScreen() *MainScreen {
	return &MainScreen{
		choices: []string{"build", "run", "shell", "services"},
		cursor:  0,
	}
}

func (m *MainScreen) Render(window *model) string {
	s := ""

	// Render the main menu
	for i, choice := range m.choices {
		indicator := " "
		style := lipgloss.NewStyle()

		// Header for docker
		if i == 0 {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("#181b21")).Render("Docker\n")
		}

		// Header for Carbon
		if i == 3 {
			s += lipgloss.NewStyle().PaddingTop(2).Foreground(lipgloss.Color("#181b21")).Render("Carbon\n")
		}

		if i == m.cursor {
			indicator = ">"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("#2391e6"))
		}

		s += fmt.Sprintf("\n%s %s", indicator, style.Render(choice))
	}

	// 30% of window width
	width := window.width * 3 / 10

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7d8899")).
		Padding(5, 2, 5, 2).
		Width(width).
		Align(lipgloss.Left)

	return style.Render(s)
}

func (m *MainScreen) Update(msg tea.Msg, original *model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Move the cursor
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ", "right", "l":
			state := ""

			// Set the state based on cursor
			switch m.cursor {
			case 0:
				state = "build"
			case 1:
				state = "run"
			case 2:
				state = "shell"
			}

			// Set the state
			original.state = state
		}
	}

	return nil
}
