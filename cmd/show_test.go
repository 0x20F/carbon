package cmd

import (
	"co2/helpers"
	"strings"
	"testing"

	"github.com/4khara/replica"
	dockerTypes "github.com/docker/docker/api/types"
)

func TestShowRunningsReturnsErrorIfNoContainersAreRunning(t *testing.T) {
	beforeCmdTest()

	// Mock the return value of the running containers
	rv := []dockerTypes.Container{}

	replica.Mocks.SetReturnValues("RunningContainers", rv)

	// Show running
	_, err := showRunning()

	// Make sure error is set
	if err == "" {
		t.Error("showRunning should return an error if no containers are running")
	}
}

func TestRunningTableHasOneRowPerContainer(t *testing.T) {
	beforeCmdTest()

	// Show running
	res, _ := showRunning()

	// Make sure there is one row per container
	if len(res.Rows()) != 5 { // 3 containers + 2 extra rows from the Header
		t.Errorf("showRunning should return one row per container but got %d", len(res.Columns[0]))
	}
}

func TestShowRunningSortsByContainerName(t *testing.T) {
	beforeCmdTest()

	// Mock the return value of the running containers
	rv := []dockerTypes.Container{
		{
			ID:    helpers.Hash("random1", 30),
			Image: "image1",
			Names: []string{"1-first-of-its-kind"},
		},
		{
			ID:    helpers.Hash("random2", 30),
			Image: "image2",
			Names: []string{"2-second-of-its-kind"},
		},
		{
			ID:    helpers.Hash("random3", 30),
			Image: "image3",
			Names: []string{"3-second-of-its-kind"},
		},
	}

	replica.Mocks.SetReturnValues("RunningContainers", rv)

	// Show running
	res, _ := showRunning()

	// Make sure the container names are sorted
	rows := res.Rows()[2:] // Skip the header

	if !strings.Contains(rows[0], "1-first-of-its-kind") {
		t.Errorf("showRunning should sort the container names. Row 1 should contain '1-first-of-its-kind' but got %s", rows[0])
	}

	if !strings.Contains(rows[1], "2-second-of-its-kind") {
		t.Error("showRunning should sort the container names")
	}

	if !strings.Contains(rows[2], "3-second-of-its-kind") {
		t.Error("showRunning should sort the container names")
	}
}
