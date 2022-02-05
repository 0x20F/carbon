package ui

import tea "github.com/charmbracelet/bubbletea"

type Screen interface {
	Render() string
	Update(tea.Msg, *model) tea.Cmd
}

func Quit(msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return true
		}
	}

	return false
}

func NavigateBack(msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			return true
		}
	}

	return false
}
