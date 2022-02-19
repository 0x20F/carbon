package cmd

import (
	"co2/database"
	"co2/helpers"
	"co2/types"
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

func TestShowStoresReturnsErrorIfNoStoresAreFound(t *testing.T) {
	// Remove all the stores from the database
	for _, store := range database.Stores() {
		database.DeleteStore(store)
	}

	// Show stores
	_, err := showStores()

	// Make sure error is set
	if err == "" {
		t.Error("showStores should return an error if no stores are found")
	}
}

func TestShowStoresShowsOneRowPerStore(t *testing.T) {
	defer afterCmdTest()

	// Add some stores to the database
	database.AddStore(types.Store{})
	database.AddStore(types.Store{})
	database.AddStore(types.Store{})

	// Show stores
	res, _ := showStores()

	// Make sure there is one row per store
	if len(res.Rows()) != 5 { // 3 stores + 2 extra rows from the Header
		t.Errorf("showStores should return one row per store but got %d", len(res.Columns[0]))
	}
}

func TestShowStoresSortsByStorePath(t *testing.T) {
	defer afterCmdTest()

	// Add some stores to the database
	database.AddStore(types.Store{Path: "store1"})
	database.AddStore(types.Store{Path: "store2"})
	database.AddStore(types.Store{Path: "store3"})

	// Show stores
	res, _ := showStores()

	// Make sure the store paths are sorted
	rows := res.Rows()[2:] // Skip the header

	if !strings.Contains(rows[0], "store1") {
		t.Errorf("showStores should sort the store paths. Row 1 should contain 'store1' but got %s", rows[0])
	}

	if !strings.Contains(rows[1], "store2") {
		t.Error("showStores should sort the store paths")
	}

	if !strings.Contains(rows[2], "store3") {
		t.Error("showStores should sort the store paths")
	}
}

func TestShowAvailableReturnsErrorIfNoCarbonServicesAreFound(t *testing.T) {
	// Remove all the stores from the database
	for _, store := range database.Stores() {
		database.DeleteStore(store)
	}

	// Show available
	_, err := showAvailable()

	// Make sure error is set
	if err == "" {
		t.Error("showAvailable should return an error if no stores are found")
	}
}

func TestShowAvailableAddsARowForEachCarbonService(t *testing.T) {
	defer afterCmdTest()

	// Mock the return value of the file system search for services
	rv := types.CarbonConfig{
		"carbon1": types.CarbonService{
			Name: "carbon1",
			Path: helpers.RandomAlphaString(50),
		},
		"carbon2": types.CarbonService{
			Name: "carbon2",
			Path: helpers.RandomAlphaString(50),
		},
		"carbon3": types.CarbonService{
			Name: "carbon3",
			Path: helpers.RandomAlphaString(50),
		},
	}

	replica.Mocks.SetReturnValues("Services", rv)

	// Show available
	res, _ := showAvailable()

	// Make sure there is one row per store
	if len(res.Rows()) != 5 { // 3 services + 2 extra rows from the Header
		t.Errorf("showAvailable should return one row per store but got %d", len(res.Columns[0]))
	}
}

func TestShowAvailableSortsByCarbonServiceName(t *testing.T) {
	beforeCmdTest()

	// Mock the return value of the file system search for services
	rv := types.CarbonConfig{
		"carbon1": types.CarbonService{
			Name: "carbon1",
			Path: helpers.RandomAlphaString(50),
		},
		"carbon2": types.CarbonService{
			Name: "carbon2",
			Path: helpers.RandomAlphaString(50),
		},
		"carbon3": types.CarbonService{
			Name: "carbon3",
			Path: helpers.RandomAlphaString(50),
		},
	}

	replica.Mocks.SetReturnValues("Services", rv)

	// Show available
	res, _ := showAvailable()

	// Make sure the Carbon service names are sorted
	rows := res.Rows()[2:] // Skip the header

	if !strings.Contains(rows[0], "carbon1") {
		t.Errorf("showAvailable should sort the Carbon service names. Row 1 should contain 'carbon1' but got %s", rows[0])
	}

	if !strings.Contains(rows[1], "carbon2") {
		t.Error("showAvailable should sort the Carbon service names")
	}

	if !strings.Contains(rows[2], "carbon3") {
		t.Error("showAvailable should sort the Carbon service names")
	}
}
