package cmd

import (
	"co2/database"
	"co2/docker"
	"co2/helpers"
	"co2/runner"
	"co2/types"
	"sync"

	"github.com/4khara/replica"
	dockerTypes "github.com/docker/docker/api/types"
)

type MockWrapperCmd struct{}

func (w *MockWrapperCmd) RunningContainers() []dockerTypes.Container {
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
			ID:    helpers.Hash("ayeeeeee lmaoooooo", 30),
			Image: "image1",
			Names: []string{"/docker-container1"},
		},
		{
			ID:    helpers.Hash("ayeeeeee lmat1r3t13ooo", 30),
			Image: "image2",
			Names: []string{"/docker-container2"},
		},
		{
			ID:    helpers.Hash("ayeeee2po3ewmföqaee lmaoooooo", 30),
			Image: "image3",
			Names: []string{"/docker-container3"},
		},
	}
}

type MockFs struct{}

func (f MockFs) Services() types.CarbonConfig {
	_, rv := replica.MockFn()

	if rv != nil {
		var carbonConfig types.CarbonConfig

		if rv[0] != nil {
			carbonConfig = rv[0].(types.CarbonConfig)
		}

		return carbonConfig
	}

	return nil
}

type MockExecutor struct{}

func (e *MockExecutor) Execute(done *sync.WaitGroup, command string, label string) {
	replica.MockFn(done, command, label)

	done.Done()
}

func beforeCmdTest() {
	WrapFs(MockFs{})

	runner.CustomExecutor(&MockExecutor{})
	docker.CustomWrapper(&MockWrapperCmd{})

	replica.Mocks.Clear()

	cleanup()
}

func afterCmdTest() {
	cleanup()
}

func cleanup() {
	// Cleanup the database
	for _, store := range database.Stores() {
		database.DeleteStore(store)
	}

	for _, container := range database.Containers() {
		database.DeleteContainer(container)
	}
}
