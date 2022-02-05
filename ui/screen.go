package ui

import tea "github.com/charmbracelet/bubbletea"

type Screen interface {
	Render(window *model) string
	Update(tea.Msg, *model) tea.Cmd
}
