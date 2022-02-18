package runner

import (
	"fmt"
	"strings"
	"sync"

	exec "github.com/go-cmd/cmd"
)

type ExecutorInterface interface {
	Execute(*sync.WaitGroup, string, string)
}

type executorImpl struct{}

func (e *executorImpl) Execute(done *sync.WaitGroup, command string, label string) {
	// Split into params
	params := strings.Split(command, " ")

	opts := exec.Options{
		Buffered:  false,
		Streaming: true,
	}
	run := exec.NewCmdOptions(opts, params[0], params[1:]...)

	// Stream output from the command and close when
	// both channels close.
	go func(doneChan *sync.WaitGroup) {
		defer doneChan.Done()

		for run.Stdout != nil || run.Stderr != nil {
			select {
			case out, ok := <-run.Stdout:
				if !ok {
					run.Stdout = nil
					continue
				}

				fmt.Println(label, string(out))
			case err, ok := <-run.Stderr:
				if !ok {
					run.Stderr = nil
					continue
				}

				fmt.Println(label, string(err))
			}
		}
	}(done)

	// Block waiting for command to exit, be stopped, or be killed
	<-run.Start()
}
