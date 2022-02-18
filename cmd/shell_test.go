package cmd

import (
	"co2/database"
	"co2/docker"
	"co2/types"
	"testing"

	"github.com/4khara/replica"
	dockerTypes "github.com/docker/docker/api/types"
)

type MockWrapperShell struct{}

func (w *MockWrapperShell) RunningContainers() []dockerTypes.Container {
	replica.MockFn()

	return []dockerTypes.Container{
		{
			ID:    "1",
			Image: "image1",
			Names: []string{"/docker-container1"},
		},
		{
			ID:    "2",
			Image: "image2",
			Names: []string{"/docker-container2"},
		},
		{
			ID:    "3",
			Image: "image3",
			Names: []string{"/docker-container3"},
		},
	}
}

func beforeShellTest() {
	docker.CustomWrapper(&MockWrapperShell{})
}

func TestGenerateCommandFindsByUid(t *testing.T) {
	beforeShellTest()

	// Get all the generated containers first so we can get the uid
	container := docker.RunningContainers()[0]

	// Generate the command with the Uid
	command := generateCommand(container.Uid, "bash")

	if command == "" {
		t.Error("generateCommand should return a command when given a Uid")
	}
}

func TestGenerateCommandFindsByName(t *testing.T) {
	beforeShellTest()

	// Get all the generated containers first so we can get the name
	container := docker.RunningContainers()[0]

	// Generate the command with the name
	command := generateCommand(container.Name, "bash")

	if command == "" {
		t.Error("generateCommand should return a command when given a name")
	}
}

func TestGenerateCommandFindsByServiceName(t *testing.T) {
	beforeShellTest()

	// Add some containers to the database
	containers := []types.Container{
		{
			Image:       "image1",
			Name:        "container1",
			Uid:         "uid1",
			ServiceName: "service1",
		},
		{
			Image:       "image2",
			Name:        "container2",
			Uid:         "uid2",
			ServiceName: "service2",
		},
	}

	for _, container := range containers {
		database.AddContainer(container)
	}

	// Generate the command with the service name
	command := generateCommand("service1", "bash")

	if command == "" {
		t.Error("generateCommand should return a command when given a service name")
	}
}

func TestGenerateCommandReturnsEmptyWhenContainerNotFound(t *testing.T) {
	beforeShellTest()

	// Generate the command with the service name
	command := generateCommand("non-existent-lol", "bash")

	if command != "" {
		t.Error("generateCommand should return an empty string when the container is not found")
	}
}

func TestByCarbonFindsByServiceName(t *testing.T) {
	beforeShellTest()

	// Add some containers to the database
	containers := []types.Container{
		{
			Image:       "image1",
			Name:        "container1",
			Uid:         "uid1",
			ServiceName: "service1",
		},
		{
			Image:       "image2",
			Name:        "container2",
			Uid:         "uid2",
			ServiceName: "service2",
		},
	}

	for _, container := range containers {
		database.AddContainer(container)
	}

	// Generate the command with the service name
	command := byCarbon("service1")

	if command == "" {
		t.Error("byCarbon should return a command when given a service name")
	}
}

func TestByDockerFindsByUid(t *testing.T) {
	beforeShellTest()

	// Get all the generated containers first so we can get the uid
	container := docker.RunningContainers()[0]

	// Generate the command with the Uid
	command := byDocker(container.Uid)

	if command == "" {
		t.Error("byDocker should return a command when given a Uid")
	}
}

func TestByDockerFindsByName(t *testing.T) {
	beforeShellTest()

	// Get all the generated containers first so we can get the name
	container := docker.RunningContainers()[0]

	// Generate the command with the name
	command := byDocker(container.Name)

	if command == "" {
		t.Error("byDocker should return a command when given a name")
	}
}
