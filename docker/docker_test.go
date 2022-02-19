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
	_, rv := replica.MockFn()

	if rv != nil {
		var containers []dockerTypes.Container

		if rv[0] != nil {
			containers = rv[0].([]dockerTypes.Container)
		}

		return containers
	}

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
	replica.Mocks.Clear()
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

func TestApiWrapperDoesntRemoveFirstContainerNameCharacterIfNoSlash(t *testing.T) {
	before()

	// Mock the return value to return one without slash
	replica.Mocks.SetReturnValues("RunningContainers", []dockerTypes.Container{
		{
			ID:    "1",
			Image: "image1",
			Names: []string{"agent1337"},
		},
	})

	containers := RunningContainers()

	// Make sure the container name isn't missing the first character
	if containers[0].Name != "agent1337" {
		t.Errorf("Expected container name to be whole even if docker didn't send back one with a slash, got %s", containers[0].Name)
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
