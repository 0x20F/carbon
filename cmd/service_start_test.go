package cmd

import (
	"co2/database"
	"co2/types"
	"testing"

	"github.com/4khara/replica"
)

func mockStores() []types.Store {
	return []types.Store{
		{
			Path: "foo",
			Env:  "FOOPATH",
		},
		{
			Path: "bar",
			Env:  "BARPATH",
		},
	}
}

func mockCarbonConfig() types.CarbonConfig {
	return types.CarbonConfig{
		"foo": types.CarbonService{
			Name: "foo",
			DependsOn: []string{
				"bar",
			},
			Store: &mockStores()[0],
			FullContents: map[string]interface{}{
				"container_name": "foo",
				"image":          "something",
			},
		},
		"bar": types.CarbonService{
			Name:  "bar",
			Store: &mockStores()[1],
			FullContents: map[string]interface{}{
				"container_name": "bar",
				"image":          "or_something",
			},
		},
		"baz": types.CarbonService{
			Name: "baz",
			DependsOn: []string{
				"bar",
			},
			Store: &mockStores()[0],
			FullContents: map[string]interface{}{
				"container_name": "baz",
				"image":          "and_something",
			},
		},
	}
}

func TestRunExecutesTheBuiltCommandWithAllServicesInOne(t *testing.T) {
	beforeCmdTest()

	envs, file, _ := compose(mockCarbonConfig())

	run(file, envs, []string{"foo", "bar", "baz"})

	if replica.Mocks.GetCallCount("Execute") != 1 {
		t.Error("run should execute the built command")
	}
}

func TestShouldRunWhenForced(t *testing.T) {
	if !shouldRun([]string{"foo", "bar"}, true) {
		t.Error("shouldRun should return true when force is true")
	}
}

func TestShouldRunWhenServicesFound(t *testing.T) {
	beforeCmdTest()

	// Add some containers to the database
	containers := []types.Container{
		{
			ServiceName: "foo",
		},
		{
			ServiceName: "bar",
		},
	}

	for _, container := range containers {
		database.AddContainer(container)
	}

	// Make sure we don't run if services are in the database
	if shouldRun([]string{"foo", "bar"}, false) {
		t.Error("shouldRun should return false when services are found")
	}
}

func TestShouldRunWhenAnyServicesNotFound(t *testing.T) {
	beforeCmdTest()

	// Add some containers to the database
	containers := []types.Container{
		{
			ServiceName: "foo",
		},
		{
			ServiceName: "bar",
		},
	}

	for _, container := range containers {
		database.AddContainer(container)
	}

	// Make sure we run if services are not in the database
	if !shouldRun([]string{"baz", "qux"}, false) {
		t.Error("shouldRun should return true when services are not found")
	}
}

func TestExtractReturnsEmptyMapWhenNoServicesAreFound(t *testing.T) {
	beforeCmdTest()

	replica.Mocks.SetReturnValues("Services", mockCarbonConfig())

	// Make sure to search for something that doesn't exist
	choices := extract([]string{"baz", "qux"})

	if len(choices) != 0 {
		t.Error("extract should return empty map when no services are found")
	}
}

func TestExtractSkipsServicesThatDependOnOtherServicesIfDependenciesAreNotPresent(t *testing.T) {
	beforeCmdTest()

	replica.Mocks.SetReturnValues("Services", mockCarbonConfig())

	// Make sure to search for something that doesn't exist
	choices := extract([]string{"foo"})

	if len(choices) != 0 {
		t.Error("extract should return empty map when services that have dependencies that are not provided are found")
	}
}

func TestExtractReturnsServicesIfDependenciesAreMet(t *testing.T) {
	beforeCmdTest()

	replica.Mocks.SetReturnValues("Services", mockCarbonConfig())

	// Make sure to search for something that doesn't exist
	choices := extract([]string{"foo", "bar"})

	if len(choices) == 0 {
		t.Error("extract should return map when dependencies are met")
	}

	if len(choices) != 2 {
		t.Error("extract should return 2 services when dependencies are met")
	}
}

func TestExtractOverridesContainerName(t *testing.T) {
	beforeCmdTest()

	replica.Mocks.SetReturnValues("Services", mockCarbonConfig())

	// Make sure to search for something that doesn't exist
	choices := extract([]string{"foo"})

	if choices["foo"].FullContents["container_name"] == "foo" {
		t.Error("extract should override the container name")
	}
}

func TestComposeReturnsErrorIfNoServicesAreFound(t *testing.T) {
	beforeCmdTest()

	_, _, err := compose(types.CarbonConfig{})
	if err == nil {
		t.Error("compose should return error when no services are found")
	}
}

func TestComposeReturnsTheRightAmountOfEnvironmentFilePaths(t *testing.T) {
	beforeCmdTest()

	paths, _, _ := compose(mockCarbonConfig())

	if len(paths) != 2 {
		t.Error("compose should return 2 paths when 2 services are found")
	}
}

func TestComposeReturnsTheRightFile(t *testing.T) {
	beforeCmdTest()

	custom := types.CarbonConfig{}
	m := mockCarbonConfig()

	// Get only the configs we care about for this test
	custom["foo"] = m["foo"]
	custom["bar"] = m["bar"]

	paths, config, _ := compose(custom)

	if len(config.Services) != 2 || len(config.Services) != len(paths) {
		t.Error("compose should return 2 env files when 2 services are found with different files are found")
	}
}

func TestComposeDoesNotDuplicateEnvironmentPathsThatAlreadyExist(t *testing.T) {
	beforeCmdTest()

	custom := types.CarbonConfig{}
	m := mockCarbonConfig()

	// Get only the configs we care about for this test
	custom["foo"] = m["foo"]
	custom["baz"] = m["baz"]

	paths, _, _ := compose(custom)

	if len(paths) != 1 {
		t.Error("compose should return 1 env files when 2 services with the same env file are found got", len(paths))
	}
}

func TestContainerizeAddsAllServicesInComposeFileAsContainers(t *testing.T) {
	beforeCmdTest()

	_, file, _ := compose(mockCarbonConfig())

	// Make sure there are no containers in the database
	if len(database.Containers()) != 0 {
		t.Error("database should be empty before containerize is called")
	}

	containerize(file)

	// Make sure there are containers in the database
	if len(database.Containers()) != len(mockCarbonConfig()) {
		t.Error("database should have containers after containerize is called")
	}
}

func TestContainerizeHashesAllContainers(t *testing.T) {
	beforeCmdTest()

	_, file, _ := compose(mockCarbonConfig())

	// Make sure there are no containers in the database
	if len(database.Containers()) != 0 {
		t.Error("database should be empty before containerize is called")
	}

	containerize(file)

	// Make sure all containers have a hash
	for _, container := range database.Containers() {
		if container.Uid == "" {
			t.Error("container should have a hash after containerize is called")
		}
	}
}

func TestContainerizeInjectsTheComposeFilePathIntoEachContainer(t *testing.T) {
	beforeCmdTest()

	_, file, _ := compose(mockCarbonConfig())

	// Make sure there are no containers in the database
	if len(database.Containers()) != 0 {
		t.Error("database should be empty before containerize is called")
	}

	containerize(file)

	// Make sure all containers have a hash
	for _, container := range database.Containers() {
		if container.ComposeFile != file.Path() {
			t.Error("container should have a hash after containerize is called")
		}
	}
}
