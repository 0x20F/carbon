package printer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Color string

// Color definitions so there's an overall
// theme that the entire application follows
// when things are being logged.
const (
	Cyan   Color = "#00ccff"
	Green  Color = "#bee38d"
	Red    Color = "#ff5252"
	Yellow Color = "#f8c76a"
	Grey   Color = "#7e8285"
)

var (
	// Header style is the first part that shows up
	// in a log message. It usually contains a single word
	// summary of what's going on and gets a background color.
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Bold(true)

	// Info style is some good to have information, however,
	// not very important to the user so it gets faded out a bit.
	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))

	// The highlight is the bit of the extra information that is
	// important to the user. Such as a number of things that the
	// info text describes.
	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f4f4f4"))

	// All the lines that are generated after the header in order
	// to show that things are happening step by step.
	extraStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#aaaaaa"))

	// The indicator color for the extra style, usually
	// an arrow.
	indicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Red))
)

// Creates an info header, allowing you to specify what color
// it should be. Could be used for Warnings, general information
// or something else.
//
// Not errors though.
func Info(color Color, title string, info string, highlight string) {
	headerStyle.Background(lipgloss.Color(color))

	out.Ln(
		headerStyle.Render(fmt.Sprintf(" %s ", title)),
		infoStyle.Render(info),
		highlightStyle.Render(highlight),
	)
}

// Creates an error header with a predefined color so that all errors
// will always look the same.
//
// Whenever you need to display an error for some reason, you should
// use this instead of the info method. It does basically the same thing
// Except it's red.
func Error(title string, info string, highlight string) {
	Info(
		Red,
		headerStyle.Render(fmt.Sprintf(" %s ", title)),
		infoStyle.Render(info),
		highlightStyle.Render(highlight),
	)
}

// Adds extra information after a header.
// All the lines printed with this will start with a custom
// indicator of your chosen color and will be indented as to
// show that they belong to the header above them.
//
// Note that the indentation doesn't use an actual tab since it's
// usually massive and looks ugly. Just 3 spaces.
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

	out.Ln(
		"  ",
		strings.Join(rendered, "\n   "),
	)
}

// Generates the styling for the text that should be printed
// but returns it instead of sending it to stdout.
//
// This allows you to fully generate styled errors and
// return them to be printed somewhere else if needed.
func Render(color Color, title, info, highlight string) string {
	headerStyle.Background(lipgloss.Color(color))

	return fmt.Sprintf(
		"%s %s %s",
		headerStyle.Render(fmt.Sprintf(" %s ", title)),
		infoStyle.Render(info),
		highlightStyle.Render(highlight),
	)
}
