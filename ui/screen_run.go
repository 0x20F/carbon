package ui

import tea "github.com/charmbracelet/bubbletea"

type RunScreen struct{}

func (m *RunScreen) Render() string {
	return "Run Menu"
}

func (m *RunScreen) Update(msg tea.Msg, original *model) tea.Cmd {
	return nil
}
