package logger

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
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

func Info(title string, info string, highlight string) {
	headerStyle.Background(lipgloss.Color("#00ccff"))

	fmt.Println(
		headerStyle.Render(fmt.Sprintf(" %s ", title)),
		infoStyle.Render(info),
		highlightStyle.Render(highlight),
	)
}
