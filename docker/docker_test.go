package docker

import (
	"co2/helpers"
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

func TestContainerKeys(t *testing.T) {
	before()

	containers := RunningContainers()

	// Make sure the keys are correct
	expected := []string{
		helpers.Hash("image1/container1", 4),
		helpers.Hash("image2/container2", 4),
	}

	found := false

	for _, key := range expected {
		if key == containers[0].Uid {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find", expected[0], "in the container keys")
	}
}
