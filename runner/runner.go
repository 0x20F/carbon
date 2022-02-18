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
func Execute(commands ...types.Command) chan struct{} {
	executor := Executor()
	done := make(chan struct{})

	for _, command := range commands {
		hash := helpers.Hash(command.Text, 14)
		color := helpers.StringToColor(hash)

		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color))

		label := ""

		if command.Name != "" {
			label = fmt.Sprintf("[ %s ]:", style.Render(command.Name))
		}

		go executor.Execute(done, command.Text, label)
	}

	return done
}
