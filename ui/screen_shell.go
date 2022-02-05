package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ShellScreen struct {
	choices []*MenuItem
	cursor  int
	input   textinput.Model
	typing  bool
}

func initialShellScreen() *ShellScreen {
	ti := textinput.New()
	ti.Placeholder = "Something..."
	ti.CharLimit = 200
	ti.Width = 20

	choices := []*MenuItem{
		NewMenuItem("bash", "Bash Shell"),
		NewMenuItem("sh", "Sh Shell"),
		NewMenuItem("zsh", "Zsh Shell"),
		NewMenuItem("custom", "Custom Shell"),
	}

	choices[0].Select()

	return &ShellScreen{
		choices: choices,
		cursor:  0,

		input: ti,
	}
}

func (m *ShellScreen) Render(window *model) string {
	// Header
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("#181b21")).Render("Shell Within Container\nPick Shell Type:\n")

	// Render the main menu
	for i, choice := range m.choices {
		if m.typing && i == 3 {
			s += "\n" + m.input.View()
			continue
		}

		s += choice.View()
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7d8899")).
		PaddingBottom(5).
		Align(lipgloss.Left)

	return style.Render(s)
}

func (m *ShellScreen) Update(msg tea.Msg, original *model) tea.Cmd {
	if m.typing {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)

		return cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Move the cursor
		case "up", "k":
			m.typing = false

			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			m.typing = false

			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ", "right", "l":
			// If custom, show text input
			if m.cursor == 3 {
				m.typing = true
				return textinput.Blink
			}
		}
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

	return nil
}
