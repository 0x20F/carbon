package docker

import (
	"co2/helpers"
	"strings"
	"testing"

	"github.com/4khara/replica"
	dockerTypes "github.com/docker/docker/api/types"
)

type MockWrapper struct{}

func (w *MockWrapper) RunningContainers() []dockerTypes.Container {
	replica.MockFn()

	return []dockerTypes.Container{
		{
			ID:    "1",
			Image: "image1",
			Names: []string{"/container1"},
		},
		{
			ID:    "2",
			Image: "image2",
			Names: []string{"/container2"},
		},
	}
}

func before() {
	CustomWrapper(&MockWrapper{})
}

func TestRunningContainers(t *testing.T) {
	before()

	containers := RunningContainers()

	// Make sure we get all the expected containers
	if len(containers) != 2 {
		t.Error("Expected 2 containers, got", len(containers))
	}
}

func TestApiWrapperRemovesContainerNameSlash(t *testing.T) {
	before()

	containers := RunningContainers()

	// Make sure none of the container names start with a slash
	for _, container := range containers {
		// Check if starts with
		if strings.HasPrefix(container.Name, "/") {
			t.Error("Expected container name to not start with a slash")
		}
	}
}

func TestContainerKeys(t *testing.T) {
	before()

	containers := RunningContainers()

	// Make sure the keys are correct
	expected := []string{
		helpers.Hash("image1container1", 4),
		helpers.Hash("image2container2", 4),
	}

	if containers[0].Uid != expected[0] {
		t.Error("Expected", expected[0], "got", containers[0].Uid)
	}

	if containers[1].Uid != expected[1] {
		t.Error("Expected", expected[1], "got", containers[1].Uid)
	}
}
