package printer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Color string

const (
	Cyan   Color = "#00ccff"
	Green  Color = "#bee38d"
	Red    Color = "#ff5252"
	Yellow Color = "#f8c76a"
	Grey   Color = "#7e8285"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f4f4f4"))

	extraStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#eaeaea"))

	indicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Red))
)

func Info(color Color, title string, info string, highlight string) {
	headerStyle.Background(lipgloss.Color(color))

	fmt.Println(
		"\n",
		headerStyle.Render(fmt.Sprintf(" %s ", title)),
		infoStyle.Render(info),
		highlightStyle.Render(highlight),
	)
}

func Error(title string, info string, highlight string) {
	headerStyle.Background(lipgloss.Color(Red))

	fmt.Println(
		"\n",
		headerStyle.Render(fmt.Sprintf(" %s ", title)),
		infoStyle.Render(info),
		highlightStyle.Render(highlight),
	)
}

func Extra(color Color, info ...string) {
	indicator.Foreground(lipgloss.Color(color))
	rendered := []string{}

	for _, str := range info {
		rendered = append(
			rendered,
			fmt.Sprintf("%s %s",
				indicator.Render("→"),
				extraStyle.Render(str),
			),
		)
	}

	fmt.Println(
		"  ",
		strings.Join(rendered, "\n   "),
	)
}
