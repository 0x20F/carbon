package runner

import (
	"co2/helpers"
	"co2/types"
	"sync"
	"testing"

	"github.com/4khara/replica"
	"github.com/charmbracelet/lipgloss"
)

type MockExecutor struct{}

func (e *MockExecutor) Execute(done *sync.WaitGroup, command string, label string) {
	replica.MockFn(done, command, label)

	done.Done()
}

func before() {
	CustomExecutor(&MockExecutor{})
}

func TestLabelGeneratesCorrectly(t *testing.T) {
	generated := label("test")

	if generated != "[ test ]:" {
		t.Errorf("Expected [ test ], got %s", generated)
	}
}

func TestColorizeAddsColorOfCommand(t *testing.T) {
	command := types.Command{
		Text: "test",
		Name: "lmao",
	}

	hashed := helpers.Hash(command.Text, 14)
	commandColor := helpers.StringToColor(hashed)
	coloredCommand := lipgloss.NewStyle().
		Foreground(lipgloss.Color(commandColor)).
		Render(command.Name)

	colored := colorize(command)

	if len(colored) != len(coloredCommand) {
		t.Errorf("Expected length %d, got %d", len(coloredCommand), len(colored))
	}
}

func TestExecuteCallsAllProvidedCommands(t *testing.T) {
	before()

	commands := []types.Command{
		{
			Text: "test",
			Name: "lmao",
		},
		{
			Text: "test2",
			Name: "lmao2",
		},
	}

	Execute(commands...)

	if replica.Mocks.GetCallCount("Execute") != len(commands) {
		t.Errorf("Expected %d calls, got %d", len(commands), replica.Mocks.GetCallCount("Execute"))
	}
}
