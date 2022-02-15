package runner

import (
	"fmt"
	"strings"

	exec "github.com/go-cmd/cmd"
)

// Executes a given command in the prefered shell
// of your platform.
//
// This will stream all the output to the console and end itself
// when the output channels have been closed.
func Execute(command string) {
	// Split into params
	params := strings.Split(command, " ")

	opts := exec.Options{
		Buffered:  false,
		Streaming: true,
	}
	run := exec.NewCmdOptions(opts, params[0], params[1:]...)

	// Stream output from the command and close when
	// both channels close.
	doneChan := make(chan struct{})
	go func(doneChan chan struct{}) {
		defer close(doneChan)

		for run.Stdout != nil || run.Stderr != nil {
			select {
			case out, ok := <-run.Stdout:
				if !ok {
					run.Stdout = nil
					continue
				}

				fmt.Println(string(out))
			case err, ok := <-run.Stderr:
				if !ok {
					run.Stderr = nil
					continue
				}

				fmt.Println(string(err))
			}
		}
	}(doneChan)

	<-run.Start()
	// Block waiting for command to exit, be stopped, or be killed
	<-doneChan
}
