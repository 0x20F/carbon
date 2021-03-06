package cmd

import (
	"co2/database"
	"co2/docker"
	"co2/types"
	"testing"
)

func TestShouldRunLogsCommandReturnsFalseWithNoCommands(t *testing.T) {
	commands := []types.Command{}

	if shouldRunLogsCommand(commands) {
		t.Error("shouldRunLogsCommand should return false when no commands are provided")
	}
}

func TestShouldRunLogsCommandReturnsTrueWithCommands(t *testing.T) {
	commands := []types.Command{
		{
			Label: "docker",
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

	if commands[0].Label != "container1" {
		t.Error("generateCommands should set the label to the container name")
	}
}

func TestContainersFilterFindsByUid(t *testing.T) {
	beforeCmdTest()

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
	beforeCmdTest()

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
	beforeCmdTest()

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
