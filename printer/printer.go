package printer

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Color string

const (
	// ANSI color codes
	Cyan  Color = "36"
	Green Color = "32"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f4f4f4"))
)

func Info(color Color, title string, info string, highlight string) {
	headerStyle.Background(lipgloss.Color(color))

	fmt.Println(
		headerStyle.Render(fmt.Sprintf(" %s ", title)),
		infoStyle.Render(info),
		highlightStyle.Render(highlight),
	)
}
