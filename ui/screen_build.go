package ui

import tea "github.com/charmbracelet/bubbletea"

type BuildScreen struct{}

func (m *BuildScreen) Render() string {
	return "Build Menu"
}

func (m *BuildScreen) Update(msg tea.Msg, original *model) tea.Cmd {
	return nil
}
