package runner

import (
	"co2/helpers"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	exec "github.com/go-cmd/cmd"
)

// Executes a given command in the prefered shell
// of your platform.
//
// This will stream all the output to the console and end itself
// when the output channels have been closed.
func singleRun(doneChan chan struct{}, command string, style lipgloss.Style) {
	// Split into params
	params := strings.Split(command, " ")

	opts := exec.Options{
		Buffered:  false,
		Streaming: true,
	}
	run := exec.NewCmdOptions(opts, params[0], params[1:]...)

	// Stream output from the command and close when
	// both channels close.
	go func(doneChan chan struct{}) {
		defer close(doneChan)

		header := fmt.Sprintf("[ %s ]", style.Render(command[:14]))

		for run.Stdout != nil || run.Stderr != nil {
			select {
			case out, ok := <-run.Stdout:
				if !ok {
					run.Stdout = nil
					continue
				}

				fmt.Println(header, string(out))
			case err, ok := <-run.Stderr:
				if !ok {
					run.Stderr = nil
					continue
				}

				fmt.Println(header, string(err))
			}
		}
	}(doneChan)

	<-run.Start()
	// Block waiting for command to exit, be stopped, or be killed
	<-doneChan
}

func Execute(commands ...string) chan struct{} {
	done := make(chan struct{})

	for _, command := range commands {
		hash := helpers.Hash(command, 14)
		color := helpers.StringToColor(hash)

		// Create a custom color based on the command to distinguish
		// the outputs
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color))

		go singleRun(done, command, style)
	}

	return done
}
