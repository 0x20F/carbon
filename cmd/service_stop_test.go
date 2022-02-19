package cmd

import (
	"co2/database"
	"co2/types"
	"testing"

	"github.com/4khara/replica"
)

func TestGroupByComposeFileUnderstandsUidsAndServiceNames(t *testing.T) {
	beforeCmdTest()

	// Add some containers to the database
	containers := []types.Container{
		{
			Image:       "image1",
			Name:        "container1",
			Uid:         "uid1",
			ServiceName: "service1",
			ComposeFile: "file1",
		},
		{
			Image:       "image2",
			Name:        "container2",
			Uid:         "uid2",
			ServiceName: "service2",
			ComposeFile: "file1",
		},
		{
			Image:       "image3",
			Name:        "container3",
			Uid:         "uid3",
			ServiceName: "service3",
			ComposeFile: "file2",
		},
	}

	for _, container := range containers {
		database.AddContainer(container)
	}

	// Get the containers grouped by the compose file
	groupedContainers := groupByComposeFile("uid1", "service3")

	// Make sure there are 2 groups
	if len(groupedContainers) != 2 {
		t.Error("groupByComposeFile should return 2 groups when given a uid and service name got", len(groupedContainers))
	}

	// Make sure there is a group for file1
	if len(groupedContainers["file1"]) != 1 {
		t.Error("groupByComposeFile should return 2 containers when given a uid and service name")
	}

	// Make sure there is a group for file2
	if len(groupedContainers["file2"]) != 1 {
		t.Error("groupByComposeFile should return 2 containers when given a uid and service name")
	}
}

func TestStopDeletesContainersFromTheDatabase(t *testing.T) {
	beforeCmdTest()

	// Add some containers to the database
	containers := []types.Container{
		{
			Image:       "image1",
			Name:        "container1",
			Uid:         "uid1",
			ServiceName: "service1",
			ComposeFile: "file1",
		},
		{
			Image:       "image2",
			Name:        "container2",
			Uid:         "uid2",
			ServiceName: "service2",
			ComposeFile: "file1",
		},
		{
			Image:       "image3",
			Name:        "container3",
			Uid:         "uid3",
			ServiceName: "service3",
			ComposeFile: "file2",
		},
	}

	for _, container := range containers {
		database.AddContainer(container)
	}

	// Get the containers grouped by the compose file
	groupedContainers := groupByComposeFile("uid1", "service3")

	// Make sure there are 2 groups
	if len(groupedContainers) != 2 {
		t.Error("groupByComposeFile should return 2 groups when given a uid and service name got", len(groupedContainers))
	}

	// Stop the containers
	stopContainers(groupedContainers)

	// Make sure the containers are gone
	if len(database.Containers()) != 1 {
		t.Error("stopContainers should remove all containers from the database")
	}

	// Make sure the compose command was executed as well, once for each compose file
	if replica.Mocks.GetCallCount("Execute") != 2 {
		t.Error("stopContainers should execute the compose command")
	}
}
