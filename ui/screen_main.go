package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var menuStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7d8899")).
	PaddingBottom(5).
	Width(WindowDimensions[0] * 3 / 10).
	Height(WindowDimensions[1]).
	Align(lipgloss.Left)

type MainScreen struct {
	choices []*MenuItem
	cursor  int
}

func initMainScreen() *MainScreen {
	choices := []*MenuItem{
		NewMenuItem("build", "Build Container"),
		NewMenuItem("run", "Run Container"),
		NewMenuItem("shell", "Shell Within Container"),
		NewMenuItem("services", "Services"),
	}

	choices[0].Select()

	return &MainScreen{
		choices: choices,
		cursor:  0,
	}
}

func (m *MainScreen) Init() tea.Cmd {
	return nil
}

func (m *MainScreen) View() string {
	// Header
	s := Title("Docker Wrappers")

	// Render the main menu
	for i, choice := range m.choices {
		if i == 3 {
			s += Title("\n\nCarbon")
		}

		s += choice.View()
	}

	return menuStyle.Render(s)
}

func (m *MainScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			State = state
		}

		// Loop through all menu items
		for i, choice := range m.choices {
			// If the cursor is on this item
			if i == m.cursor {
				choice.Select()
			} else {
				choice.Deselect()
			}

			choice.Update(msg)
		}
	}

	return m, nil
}
