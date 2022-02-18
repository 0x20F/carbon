package cmd

import (
	"co2/database"
	"co2/docker"
	"co2/types"
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

func before() {
	docker.CustomWrapper(&MockWrapper{})
}

func TestShouldRunLogsCommandReturnsFalseWithNoCommands(t *testing.T) {
	commands := []types.Command{}

	if shouldRunLogsCommand(commands) {
		t.Error("shouldRunLogsCommand should return false when no commands are provided")
	}
}

func TestShouldRunLogsCommandReturnsTrueWithCommands(t *testing.T) {
	commands := []types.Command{
		{
			Name: "docker",
		},
	}

	if !shouldRunLogsCommand(commands) {
		t.Error("shouldRunLogsCommand should return true when commands are provided")
	}
}

func TestGenerateCommandsReturnsSameAmountOfCommandsAsContainers(t *testing.T) {
	containers := []types.Container{
		{
			Name: "container1",
		},
		{
			Name: "container2",
		},
	}

	commands := generateCommands(containers, false)

	if len(commands) != len(containers) {
		t.Error("generateCommands should return the same amount of commands as containers")
	}
}

func TestCommandLabelMatchesContainerName(t *testing.T) {
	containers := []types.Container{
		{
			Name: "container1",
		},
		{
			Name: "container2",
		},
	}

	commands := generateCommands(containers, true)

	if commands[0].Name != "container1" {
		t.Error("generateCommands should set the label to the container name")
	}
}

func TestContainersFilterFindsByUid(t *testing.T) {
	before()

	// Get the containers that docker will return so we can get the hashes
	dockerContainers := docker.RunningContainers()

	// Filter the containers
	choices := []string{dockerContainers[1].Uid, dockerContainers[0].Uid}
	filtered := filterContainers(choices)

	if len(filtered) != 2 {
		t.Error("filterContainers should return 2 containers")
	}
}

func TestContainersFilterFindsByCarbonServiceName(t *testing.T) {
	before()

	containers := []types.Container{
		{
			ServiceName: "container1",
			Uid:         "uid1",
		},
		{
			ServiceName: "container2",
			Uid:         "uid2",
		},
	}

	// Add some containers to the database
	for _, container := range containers {
		database.AddContainer(container)
	}

	filtered := filterContainers([]string{"container1", "container2"})

	if len(filtered) != 2 {
		t.Error("filterContainers should return 2 containers")
	}
}

func TestContainersMatchesBothUidAndServiceName(t *testing.T) {
	before()

	containers := []types.Container{
		{
			ServiceName: "container1",
			Uid:         "uid1",
		},
		{
			ServiceName: "container2",
			Uid:         "uid2",
		},
	}

	// Add some containers to the database
	for _, container := range containers {
		database.AddContainer(container)
	}

	// Get the containers that docker will return so we can get the hashes
	dockerContainers := docker.RunningContainers()

	// Filter the containers
	choices := []string{dockerContainers[1].Uid, "container1"}
	filtered := filterContainers(choices)

	if len(filtered) != 2 {
		t.Error("filterContainers should return 2 containers")
	}
}
