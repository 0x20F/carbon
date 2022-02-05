package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var itemSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2391e6"))
var itemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7d8899"))

type MenuItem struct {
	name        string
	description string
	selected    bool
}

func (m *MenuItem) Select()   { m.selected = true }
func (m *MenuItem) Deselect() { m.selected = false }

func NewMenuItem(name, description string) *MenuItem {
	return &MenuItem{
		name:        name,
		description: description,
	}
}

func (m *MenuItem) Init() tea.Cmd {
	return nil
}

func (m *MenuItem) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *MenuItem) View() string {
	if m.selected {
		return "\n> " + itemSelectedStyle.Render(m.name)
	}

	return itemStyle.Render("\n  " + m.name)
}
