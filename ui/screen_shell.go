package ui

import tea "github.com/charmbracelet/bubbletea"

type ShellScreen struct{}

func (m *ShellScreen) Render() string {
	return "Run Menu"
}

func (m *ShellScreen) Update(msg tea.Msg, original *model) tea.Cmd {
	return nil
}
