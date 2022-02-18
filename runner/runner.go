package runner

import (
	"co2/helpers"
	"co2/types"
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

var once sync.Once
var instance ExecutorInterface

// Creates a new instance of the default executors or
// returns an already created instance.
//
// Use this whenever you need to run a command somewhere.
func Executor() ExecutorInterface {
	if instance != nil {
		return instance
	}

	e := &executorImpl{}

	return CustomExecutor(e)
}

// Creates a new executor with a custom implementation.
//
// Never use this, unless writing tests.
// This is what the Executor() function runs anyway so it's better
// to let it inject all the required things.
func CustomExecutor(e ExecutorInterface) ExecutorInterface {
	once.Do(func() {
		instance = e
	})

	return instance
}

// Executes the given commands in the prefered shell
// of your platform.
//
// This will stream all the output to the console and end itself
// when the output channels have been closed.
func Execute(commands ...types.Command) {
	executor := Executor()
	var wg sync.WaitGroup

	for _, command := range commands {
		wg.Add(1)

		colored := colorize(command)
		label := label(colored)

		go executor.Execute(&wg, command.Text, label)
	}

	wg.Wait()
}

// Return a colored string that contains the given command name.
// The string is colored based on the command Text/contents meaning if
// the command is always the same, the color will also always
// be the same.
//
// To be used in combination with the label() function.
func colorize(command types.Command) string {
	hash := helpers.Hash(command.Text, 14)
	color := helpers.StringToColor(hash)
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color))

	return style.Render(command.Name)
}

// Cretes a label from a given string,
// to be shown before each line of output for a specific
// command.
func label(from string) string {
	if from == "" {
		return ""
	}

	return fmt.Sprintf("[ %s ]:", from)
}
