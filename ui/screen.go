package ui

import tea "github.com/charmbracelet/bubbletea"

type Screen interface {
	Render() string
	Update(tea.Msg, *model) tea.Cmd
}
