package cmd

import (
	"co2/database"
	"co2/docker"
	"co2/types"
	"testing"
)

func TestGenerateCommandFindsByUid(t *testing.T) {
	beforeCmdTest()

	// Get all the generated containers first so we can get the uid
	container := docker.RunningContainers()[0]

	// Generate the command with the Uid
	command := generateShellCommand(container.Uid, "bash")

	if command == "" {
		t.Error("generateCommand should return a command when given a Uid")
	}
}

func TestGenerateCommandFindsByName(t *testing.T) {
	beforeCmdTest()

	// Get all the generated containers first so we can get the name
	container := docker.RunningContainers()[0]

	// Generate the command with the name
	command := generateShellCommand(container.Name, "bash")

	if command == "" {
		t.Error("generateCommand should return a command when given a name")
	}
}

func TestGenerateCommandFindsByServiceName(t *testing.T) {
	beforeCmdTest()

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
	command := generateShellCommand("service1", "bash")

	if command == "" {
		t.Error("generateCommand should return a command when given a service name")
	}
}

func TestGenerateCommandReturnsEmptyWhenContainerNotFound(t *testing.T) {
	beforeCmdTest()

	// Generate the command with the service name
	command := generateShellCommand("non-existent-lol", "bash")

	if command != "" {
		t.Error("generateCommand should return an empty string when the container is not found")
	}
}

func TestByCarbonFindsByServiceName(t *testing.T) {
	beforeCmdTest()

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
	beforeCmdTest()

	// Get all the generated containers first so we can get the uid
	container := docker.RunningContainers()[0]

	// Generate the command with the Uid
	command := byDocker(container.Uid)

	if command == "" {
		t.Error("byDocker should return a command when given a Uid")
	}
}

func TestByDockerFindsByName(t *testing.T) {
	beforeCmdTest()

	// Get all the generated containers first so we can get the name
	container := docker.RunningContainers()[0]

	// Generate the command with the name
	command := byDocker(container.Name)

	if command == "" {
		t.Error("byDocker should return a command when given a name")
	}
}
