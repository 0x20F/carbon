package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MainScreen struct {
	choices []string
	cursor  int
}

func initMainScreen() *MainScreen {
	return &MainScreen{
		choices: []string{"build", "run", "shell"},
		cursor:  0,
	}
}

func (m *MainScreen) Render() string {
	s := ""

	// Render the main menu
	for i, choice := range m.choices {
		indicator := " "
		if i == m.cursor {
			indicator = ">"
		}

		s += fmt.Sprintf("%s %s\n", indicator, choice)
	}

	return s
}

func (m *MainScreen) Update(msg tea.Msg, original *model) tea.Cmd {
	switch msg := msg.(type) {
	// Is it a keypress?
	case tea.KeyMsg:
		// What key tho?
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

		case "enter", " ":
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
